package httprerouter

import "regexp"
import "log"
import "net/http"

type Handler func(http.ResponseWriter, *http.Request, Params)

type Route struct {
  method string
  pattern string
  re *regexp.Regexp
  cb Handler
}

type Router struct {
  routes []Route
}

type Param struct {
  Name string
  Value string
}

type Params struct {
  values []Param
}

func NewParams() *Params {
  ps := &Params{}
  ps.values = make([]Param, 0)
  return ps
}

func (ps *Params) AddParam(Name string, Value string) {
  ps.values = append(ps.values, Param{Name: Name, Value: Value})
}

func (ps *Params) ByName(Name string) *string {
  var result *string = nil

  for _, p := range ps.values {
      if p.Name == Name {
          result = &p.Value
      }
  }

  return result
}

func NewRouter() *Router {
  r := &Router{}
  r.routes = make([]Route, 0)
  return r
}

func (r *Router) Handle(method string, pattern string, cb Handler) {

  re, err := regexp.Compile(pattern)
  if err != nil {
    log.Fatal(err)
  }

  route := Route{
    method: method,
    pattern: pattern,
    re: re,
    cb: cb,
  }

  r.routes = append(r.routes, route)

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  var complete bool = false

  for _, route := range r.routes {
    log.Printf("route = %v", route)

    if req.Method == route.method {
      path := req.URL.Path

      match := route.re.FindAllStringSubmatch(path, -1)
      if len(match) > 0 {
        names := route.re.SubexpNames()
        groups := match[0]

        log.Printf("names = %v, groups = %v", names, groups)

        p := Params{}

        for i, n := range groups {
            p_name := names[i]
            p.AddParam(p_name, n)
        }

        route.cb(w, req, p)

        complete = true
      }
    }


  }

  /* No match */
  if !complete {
    http.Error(w, "", http.StatusNotFound)
  }

}
