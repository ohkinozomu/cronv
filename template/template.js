function Trim(v) {
  return v.trim();
}

function IsRunningEveryMinutes(cron) {
  // cron.split(' ').forEach((v, i, a) => {
  //     if (v != '*' && (i > 0 || v != '*/1')) {
	//   		return false;
	//   	}
  //   }
  // )
  // return true;
}

function CronvIter(cronv) {
  // TODO
}

function DateFormat() {
  // TODO
}

function NewJsDate() {
  // TODO
}

google.charts.load("current", {packages:["timeline"]});
google.charts.setOnLoadCallback(function() {
  var container = document.getElementById('cronv-timeline');
  var chart = new google.visualization.Timeline(container);
  var dataTable = new google.visualization.DataTable();
    dataTable.addColumn({ type: 'string', id: 'job' });
    dataTable.addColumn({ type: 'string', id: 'dummy bar label' });
    dataTable.addColumn({ type: 'string', role: 'tooltip' });
    dataTable.addColumn({ type: 'date', id: 'Start' });
    dataTable.addColumn({ type: 'date', id: 'End' });

  var tasks = {};
  {{ $timeFrom := .TimeFrom }}
  {{ $timeTo := .TimeTo }}
  {{range $index, $cronv := .CronEntries}}
    job = Trim('{{$cronv.Crontab.Job}}');
    tasks[job] = tasks[job] || [];
    {{if IsRunningEveryMinutes $cronv.Crontab }}
      tasks[job].push([job, '', `Every minutes ${job}`, {{NewJsDate $timeFrom}}, {{NewJsDate $timeTo}}]);
    {{else}}
      {{range CronvIter $cronv}}tasks[job].push([job, '', `{{DateFormat .Start "15:04"}} ${job}`, {{NewJsDate .Start}}, {{NewJsDate .End}}]);{{end}}
    {{ end }}
  {{end}}

  var taskByJobCount = [];
  for (var k in tasks) taskByJobCount.push({name: k, size: tasks[k].length});
  taskByJobCount.sort(function(a, b) {
    if (a.size == b.size) return 0;
    return a.size > b.size ? -1 : 1;
  });

  var rows = [];
  for (var i = 0; i < taskByJobCount.length; i++) {
    jobs = tasks[taskByJobCount[i].name];
    var jl = jobs.length;
    for (var j = 0; j < jl; j++) rows.push(jobs[j]);
  }

  if (rows.length > 0) {
    dataTable.addRows(rows);
    chart.draw(dataTable, {
      timeline: {
        colorByRowLabel: true
      },
      avoidOverlappingGridLines: false
    });
  } else {
    container.innerHTML = '<div class="alert alert-success"><strong>Woops!</strong> There is no data!</div>';
  }

  var mousePosX = undefined,
    mousePosY = undefined;

  google.visualization.events.addListener(chart, 'onmouseover', function(e) {
    var t = document.getElementsByClassName("google-visualization-tooltip")[0];
    if (mousePosX) t.style.left = mousePosX + 'px';
    if (mousePosY) t.style.top = mousePosY - 120 + 'px';
  });

  document.addEventListener('mousemove', function(e) {
    mousePosX = e.pageX;
    mousePosY = e.pageY;
  });
});