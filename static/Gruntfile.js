module.exports = function(grunt) {
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n'
            },
            build: {
                src: 'build/<%= pkg.name %>.js',
                dest: 'build/<%= pkg.name %>.min.js'
            }
        },
        jst: {
            compile: {
                files: {
                    "templates/build.js": ["templates/*.html"]
                }
            }
        },
        concat: {
            dist: {
                src: [
                    'lib/*.js',
                    'templates/*.js',
                    'src/*.js'
                ],
                dest: 'build/<%= pkg.name %>.js',
            }
        },
    });

    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-jst');

    grunt.registerTask('default', ['jst', 'concat', 'uglify']);
};
