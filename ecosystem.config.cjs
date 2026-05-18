module.exports = {
    apps: [
      {
        name: "remotemouse-server",
        script: "./apps/desktop/cmd.exe",
        watch: false,
        restart_delay: 3000,
      },
      {
        name: "remotemouse-mobile",
        script: "./node_modules/vite/bin/vite.js",
        args: "dev --host",
        cwd: "./apps/mobile/",
        watch: false,
      }
    ]
  }