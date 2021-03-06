var gulp = require('gulp');
var child = require('child_process');
var gulpsync = require('gulp-sync')(gulp);
var gutil = require('gulp-util');

var server = null;

gulp.task('default', ['assets:watch', 'server:watch']);
gulp.task('build', ['assets:copy', 'server:build']);

/*==============================
=            ASSETS            =
==============================*/

gulp.task('assets:watch', function() {
  gulp.watch('./js_src/**/*', ['assets:copy'])
})

gulp.task('assets:copy', function() {
  return gulp.src('./js_src/**/*')
    .pipe(gulp.dest('./public'))
})



/*==============================
=            SERVER            =
==============================*/

gulp.task('server:build', function() {
  var build = child.spawnSync('go', ['install']);

  if (build.stderr.length) {
    gutil.log('go build error');
  }

  return build;
});

gulp.task('server:spawn', function() {
  if (server) {
    server.kill();
  }

  server = child.spawn('./bin/deluge-server');

  server.stdout.on('data', function(data) {
    var lines = data.toString().split('\n');
    for (var i=0; i<lines.length; i++) {
      if (lines[i].length) {
        gutil.log(lines[i]);
      }
    }
  });

  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  })
});

gulp.task('server:watch', function() {
  gulp.watch([
    './*.json',
  ], ['server:spawn']);

  gulp.watch([
    './*.go',
  ], gulpsync.sync([
    'server:build',
    'server:spawn',
  ]));
});