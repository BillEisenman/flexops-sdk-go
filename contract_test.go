package flexops_test

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	flexops "github.com/BillEisenman/flexops-sdk-go/v3"
)

// testdata/gateway-operations.json is the sorted "METHOD /path" list of the Developer Docs
// OpenAPI snapshot (FlexOpsDeveloperDocs static/openapi/gateway.json, commit f6c9808).
// Refresh it by regenerating that list from the current snapshot.

// deprecatedProxyWrites are the only methods allowed to send a non-GET request through
// the legacy proxy: the carrier writes restored from 2.x, deprecated and user-token only.
var deprecatedProxyWrites = []string{
	"Carriers.USPS.GetDomesticRates", "Carriers.USPS.GetDomesticProducts", "Carriers.USPS.GetDomesticPrices",
	"Carriers.USPS.GetInternationalRates", "Carriers.USPS.GetInternationalPrices", "Carriers.USPS.CancelDomesticLabel",
	"Carriers.USPS.CancelInternationalLabel", "Carriers.USPS.CreatePickup", "Carriers.USPS.CancelPickup",
	"Carriers.USPS.CreateScanForm",
	"Carriers.UPS.ValidateAddress", "Carriers.UPS.GetRates", "Carriers.UPS.CreatePickup", "Carriers.UPS.CancelPickup",
	"Carriers.UPS.GetTransitTimes", "Carriers.UPS.GetLandedCost", "Carriers.UPS.SearchLocations",
	"Carriers.UPS.UploadDocument", "Carriers.UPS.CreateFreightShipment", "Carriers.UPS.GetFreightRate",
	"Carriers.DHL.GetMultiPieceRates", "Carriers.DHL.CreatePickup", "Carriers.DHL.CancelPickup",
	"Carriers.DHL.CalculateLandedCost", "Carriers.DHL.ScreenShipment", "Carriers.DHL.UploadInvoice",
	"Carriers.FedEx.ValidateAddress", "Carriers.FedEx.ValidatePostalCode", "Carriers.FedEx.GetRates",
	"Carriers.FedEx.CancelShipment", "Carriers.FedEx.ValidateShipment", "Carriers.FedEx.Track",
	"Carriers.FedEx.TrackMultiPiece", "Carriers.FedEx.RegisterTrackingNotification", "Carriers.FedEx.CreatePickup",
	"Carriers.FedEx.CancelPickup", "Carriers.FedEx.SearchLocations", "Carriers.FedEx.GetServiceStandards",
	"Carriers.FedEx.GetFreightRate", "Carriers.FedEx.CreateFreightShipment", "Carriers.FedEx.GroundClose",
	"Carriers.FedEx.UploadTradeDocuments", "Carriers.FedEx.CreateOpenShipment",
	"Carriers.FedEx.AddPackagesToOpenShipment", "Carriers.FedEx.ConfirmOpenShipment",
}

// closedLabelSuffixes are the proxy label-generation routes the Gateway closes with 410
// in every environment (ApiProxyLegacyRouteSemantics; matched case-insensitively).
var closedLabelSuffixes = []string{
	"generateNewUpsShipLabel", "postNewShipLabelOrder", "postUspsGenerateDomesticShippingLabel",
	"postUspsGenerateDomesticReturnsShippingLabel", "postUspsGenerateInternationalShippingLabel",
	"postUspsGenerateDomesticIndiciaLabel", "postFedExCreateNewShipment", "postFedExCreateNewReturnShipment",
	"postDhlCreateShipment", "postOnTracCreateShipment", "postLsoCreateShipment", "postAusPostCreateShipment",
	"postCanadaPostCreateShipment", "postRoyalMailCreateShipment",
}

// deprecatedMethods returns "ReceiverType.Method" for every method in the package source
// whose doc comment has a "Deprecated:" paragraph.
func deprecatedMethods(t *testing.T) map[string]bool {
	t.Helper()
	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	deprecated := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || fn.Doc == nil {
					continue
				}
				text := fn.Doc.Text()
				if !strings.HasPrefix(text, "Deprecated: ") && !strings.Contains(text, "\nDeprecated: ") {
					continue
				}
				recv := fn.Recv.List[0].Type
				if star, ok := recv.(*ast.StarExpr); ok {
					recv = star.X
				}
				if ident, ok := recv.(*ast.Ident); ok {
					deprecated[ident.Name+"."+fn.Name.Name] = true
				}
			}
		}
	}
	return deprecated
}

type documentedOperation struct {
	method  string
	pattern *regexp.Regexp
}

func loadDocumentedOperations(t *testing.T) []documentedOperation {
	t.Helper()
	data, err := os.ReadFile("testdata/gateway-operations.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Operations []string `json:"operations"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	param := regexp.MustCompile(`\{[^}]+\}`)
	ops := make([]documentedOperation, 0, len(fixture.Operations))
	for _, op := range fixture.Operations {
		method, path, _ := strings.Cut(op, " ")
		path = normalizeContractPath(path)
		var pattern string
		if path == "/api/apiproxy/{path}" {
			// The legacy proxy forwards any path below it.
			pattern = `^/api/apiproxy/.+$`
		} else {
			parts := param.Split(path, -1)
			for i := range parts {
				parts[i] = regexp.QuoteMeta(parts[i])
			}
			// Each {param} matches exactly one concrete segment.
			pattern = "^" + strings.Join(parts, `[^/]+`) + "$"
		}
		ops = append(ops, documentedOperation{method, regexp.MustCompile(pattern)})
	}
	return ops
}

// normalizeContractPath drops the query, maps /api/v1 to /api and lower-cases the path
// (ASP.NET routing is case-insensitive and the snapshot mixes casings).
func normalizeContractPath(path string) string {
	path, _, _ = strings.Cut(path, "?")
	path = strings.ToLower(path)
	if strings.HasPrefix(path, "/api/v1/") {
		path = "/api/" + strings.TrimPrefix(path, "/api/v1/")
	}
	return path
}

// contractArgs builds dummy arguments for a method: a background context, "1" for
// strings (and one value for a variadic ...string), 1 for integers and zero values for
// everything else.
func contractArgs(t *testing.T, name string, m reflect.Value) ([]reflect.Value, bool) {
	mt := m.Type()
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	args := make([]reflect.Value, mt.NumIn())
	for i := 0; i < mt.NumIn(); i++ {
		in := mt.In(i)
		switch {
		case in == ctxType:
			args[i] = reflect.ValueOf(context.Background())
		case mt.IsVariadic() && i == mt.NumIn()-1:
			if in.Elem().Kind() != reflect.String {
				t.Fatalf("%s: no dummy for variadic %s; extend contractArgs", name, in)
			}
			args[i] = reflect.ValueOf([]string{"1"}).Convert(in)
		case in.Kind() == reflect.String:
			args[i] = reflect.ValueOf("1").Convert(in)
		case in.Kind() >= reflect.Int && in.Kind() <= reflect.Int64:
			args[i] = reflect.ValueOf(int64(1)).Convert(in)
		case in.Kind() == reflect.Bool, in.Kind() == reflect.Struct, in.Kind() == reflect.Slice,
			in.Kind() == reflect.Map, in.Kind() == reflect.Pointer, in.Kind() == reflect.Interface:
			args[i] = reflect.Zero(in)
		default:
			t.Fatalf("%s: no dummy for parameter %s; extend contractArgs", name, in)
		}
	}
	return args, mt.IsVariadic()
}

// contractServices returns every service reachable from the client's exported fields,
// including nested services such as Carriers.USPS, keyed by their field path.
func contractServices(client *flexops.Client) map[string]reflect.Value {
	services := map[string]reflect.Value{}
	var walk func(prefix string, v reflect.Value)
	walk = func(prefix string, v reflect.Value) {
		s := v.Elem()
		for i := 0; i < s.NumField(); i++ {
			field := s.Type().Field(i)
			fv := s.Field(i)
			if !field.IsExported() || fv.Kind() != reflect.Pointer || fv.Elem().Kind() != reflect.Struct {
				continue
			}
			name := strings.TrimPrefix(prefix+"."+field.Name, ".")
			services[name] = fv
			walk(name, fv)
		}
	}
	walk("", reflect.ValueOf(client))
	return services
}

// Every public SDK method must call an operation in the Developer Docs OpenAPI snapshot.
// Only the deprecatedProxyWrites (each marked Deprecated in its doc comment) may write
// through the legacy proxy, which refuses API-key writes; none may send a PATCH there or
// call a 410-closed label route. Methods are found by reflection, so a new method is
// exercised automatically; one that sends no request, or takes an argument the test
// cannot fake, fails the test.
func TestContract_EveryPublicMethodCallsADocumentedOperation(t *testing.T) {
	ops := loadDocumentedOperations(t)

	var mu sync.Mutex
	current := ""
	seen := map[string][]string{}
	client, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen[current] = append(seen[current], r.Method+" "+r.URL.Path)
		mu.Unlock()
		writeJSON(w, 200, map[string]any{})
	}))
	flexops.SetRetryBaseDelay(client, time.Millisecond)

	services := contractServices(client)
	allowed := map[string]bool{}
	for _, name := range deprecatedProxyWrites {
		allowed[name] = true
	}
	proxyWriters := map[string]bool{}
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)

	methods, checked := 0, 0
	for _, serviceName := range names {
		service := services[serviceName]
		for i := 0; i < service.NumMethod(); i++ {
			name := serviceName + "." + service.Type().Method(i).Name
			m := service.Method(i)
			args, variadic := contractArgs(t, name, m)
			mu.Lock()
			current = name
			mu.Unlock()
			if variadic {
				m.CallSlice(args)
			} else {
				m.Call(args)
			}
			methods++

			mu.Lock()
			calls := seen[name]
			mu.Unlock()
			if len(calls) == 0 {
				t.Errorf("%s sent no request; give it dummy arguments that reach the API", name)
			}
			for _, call := range calls {
				checked++
				method, path, _ := strings.Cut(call, " ")
				normalized := normalizeContractPath(path)
				if strings.HasPrefix(normalized, "/api/apiproxy/") && method != http.MethodGet {
					proxyWriters[name] = true
					if !allowed[name] {
						t.Errorf("%s writes through the legacy proxy: %s %s", name, method, path)
					}
					if method == http.MethodPatch {
						t.Errorf("%s sends a PATCH through the legacy proxy, which has no PATCH route", name)
					}
					for _, suffix := range closedLabelSuffixes {
						if strings.HasSuffix(normalized, strings.ToLower(suffix)) {
							t.Errorf("%s calls the closed (410) label route %s", name, path)
						}
					}
				}
				documented := false
				for _, op := range ops {
					if op.method == method && op.pattern.MatchString(normalized) {
						documented = true
						break
					}
				}
				if !documented {
					t.Errorf("%s calls %s %s, which the Gateway OpenAPI does not document", name, method, path)
				}
			}
		}
	}
	// The allowlist is exact: each entry writes through the proxy and is marked Deprecated.
	deprecated := deprecatedMethods(t)
	for _, name := range deprecatedProxyWrites {
		if !proxyWriters[name] {
			t.Errorf("%s is allowlisted but sends no proxy write; remove it from deprecatedProxyWrites", name)
		}
		dot := strings.LastIndex(name, ".")
		service, ok := services[name[:dot]]
		if !ok || !deprecated[service.Type().Elem().Name()+"."+name[dot+1:]] {
			t.Errorf("%s writes through the legacy proxy but is not marked Deprecated", name)
		}
	}
	if methods < 50 {
		t.Fatalf("found only %d public methods; the service walk is broken", methods)
	}
	t.Logf("checked %d calls from %d public methods on %d services against %d documented operations", checked, methods, len(services), len(ops))
}
