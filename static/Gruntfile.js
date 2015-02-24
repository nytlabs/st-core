module.exports = function(grunt) {
    grunt.initConfig({
        react: {
            jsx: {
                files: [{
                    expand: true,
                    cwd: 'js',
                    src: ['**/*.jsx'],
                    dest: 'js',
                    ext: '.js'
                }]
            }
        },
    });

    grunt.loadNpmTasks('grunt-react');

    grunt.registerTask('default', ['react']);
};
