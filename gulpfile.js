var gulp = require('gulp'),
    postcss = require('gulp-postcss'),
    concatCss = require('gulp-concat-css'),
    cssnano = require('gulp-cssnano');
    babel = require('gulp-babel');
    minifyjs = require('gulp-js-minify');

gulp.task('dev-css', function () {
    return gulp.src('./src/css/page.css')
      .pipe(postcss([
        require('tailwindcss'),
        require('autoprefixer'),
      ]))
      .pipe(concatCss('page.css'))
      .pipe(cssnano({
        reduceIdents: false,
        discardComments: {removeAll: true}
      }))
      .pipe(gulp.dest('./static/css/'))
});

gulp.task('build-css', function () {
    return gulp.src('./src/css/page.css')
      .pipe(postcss([
        require('tailwindcss'),
        require('autoprefixer'),
      ]))
      .pipe(concatCss('page.css'))
      .pipe(cssnano({
        reduceIdents: false,
        discardComments: {removeAll: true}
      }))
      .pipe(gulp.dest('./static/css/'))
});

gulp.task('js', function () {
  return gulp.src("src/js/*.js")
    .pipe(babel())
    .pipe(minifyjs())
    .pipe(gulp.dest('./static/js/'));
});

gulp.task('watch-css', function() {
  gulp.watch('./src/css/*.css', gulp.series('dev-css'));
});

gulp.task('dev', gulp.series('dev-css', 'js'));

gulp.task('build', gulp.series('build-css', 'js'));