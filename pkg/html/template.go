package html

import (
	"html/template"
)

var templateHTML *template.Template = template.Must(template.New("html").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/spiermar/d3-flame-graph@2.0.3/dist/d3-flamegraph.css">

    <style>
      body {
        padding-top: 20px;
        padding-bottom: 20px;
      }

      .header {
        padding-bottom: 20px;
        padding-right: 15px;
        padding-left: 15px;
        border-bottom: 1px solid #e5e5e5;
      }

      .header .logo {
        text-decoration: none;
      }

      .header h3 {
        margin-top: 0;
        margin-bottom: 0;
        line-height: 40px;
      }

      .container {
        max-width: 990px;
      }

      table {
        font-family: monospace;
      }

      table th {
        width: 250px;
      }

      span {
        font-weight: 700;
        font-family: monospace;
      }
    </style>

    <title>pg_flame</title>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
    <div class="container">
      <div class="header clearfix">
        <nav>
          <div class="pull-right">
            <form class="form-inline" id="form">
              <a class="btn" href="javascript: resetZoom();">Reset zoom</a>
              <a class="btn" href="javascript: clear();">Clear</a>
              <div class="form-group">
                <input type="text" class="form-control" id="term">
              </div>
              <a class="btn btn-primary" href="javascript: search();">Search</a>
            </form>
          </div>
        </nav>
        <a class="logo" href="https://github.com/mgartner/pg_flame">
          <h3 class="text-muted">pg_flame</h3>
        </a>
      </div>
      <div id="chart">
      </div>
      <hr>
      <div id="details">
      </div>
    </div>

    <script src="https://d3js.org/d3.v4.min.js" charset="utf-8"></script>
    <script type="text/javascript" src=https://cdnjs.cloudflare.com/ajax/libs/d3-tip/0.9.1/d3-tip.min.js></script>
    <script type="text/javascript" src="https://cdn.jsdelivr.net/gh/spiermar/d3-flame-graph@2.1.3/dist/d3-flamegraph.min.js"></script>

    <script type="text/javascript">
    var flameGraph = d3.flamegraph()
      .width(960)
      .cellHeight(18)
      .transitionDuration(750)
      .minFrameSize(5)
      .transitionEase(d3.easeCubic)
      .sort(false)
      .title("")
      .differential(false)
      .selfValue(false)
      .setColorMapper(function(d, originalColor) {
        return d.data.color || originalColor;
      });

    var tip = d3.tip()
      .direction("s")
      .offset([8, 0])
      .attr('class', 'd3-flame-graph-tip')
      .html(function(d) {
        return d.data.name + " | " + d.data.time + "ms";
      });
    flameGraph.tooltip(tip);

    var details = document.getElementById("details");
    flameGraph.setDetailsElement(details);

    var label = function(d) {
      return d.data.detail;
    }
    flameGraph.label(label);

    var data = {{.}};

    d3.select("#chart")
      .datum(data)
      .call(flameGraph);

    document.getElementById("form").addEventListener("submit", function(event){
      event.preventDefault();
      search();
    });

    function search() {
      var term = document.getElementById("term").value;
      flameGraph.search(term);
    }

    function clear() {
      document.getElementById('term').value = '';
      flameGraph.clear();
    }

    function resetZoom() {
      flameGraph.resetZoom();
    }
    </script>
  </body>
</html>
`))

var templateTable *template.Template = template.Must(template.New("table").Parse(`
<table class="table table-striped table-bordered">
  <tbody>
    {{if .Filter}}
      <tr>
        <th>Filter</th>
        <td>{{.Filter}}</td>
      </tr>
    {{end}}
    {{if .ParentRelationship}}
      <tr>
        <th>Parent Relationship</th>
        <td>{{.ParentRelationship}}</td>
      </tr>
    {{end}}
    {{if .JoinFilter}}
      <tr>
        <th>Join Filter</th>
        <td>{{.JoinFilter}}</td>
      </tr>
    {{end}}
    {{if .HashCond}}
      <tr>
        <th>Hash Cond</th>
        <td>{{.HashCond}}</td>
      </tr>
    {{end}}
    {{if .IndexCond}}
      <tr>
        <th>Index Cond</th>
        <td>{{.IndexCond}}</td>
      </tr>
    {{end}}
    {{if .RecheckCond}}
      <tr>
        <th>Recheck Cond</th>
        <td>{{.RecheckCond}}</td>
      </tr>
    {{end}}
    {{if .BuffersHit}}
      <tr>
        <th>Buffers Shared Hit</th>
        <td>{{.BuffersHit}}</td>
      </tr>
    {{end}}
    {{if .BuffersRead}}
      <tr>
        <th>Buffers Shared Read</th>
        <td>{{.BuffersRead}}</td>
      </tr>
    {{end}}
    {{if .HashBuckets}}
      <tr>
        <th>Hash Buckets</th>
        <td>{{.HashBuckets}}</td>
      </tr>
    {{end}}
    {{if .HashBatches}}
      <tr>
        <th>Hash Batches</th>
        <td>{{.HashBatches}}</td>
      </tr>
    {{end}}
    {{if .MemoryUsage}}
      <tr>
        <th>Memory Usage</th>
        <td>{{.MemoryUsage}}kB</td>
      </tr>
    {{end}}
  </tbody>
</table>
`))
