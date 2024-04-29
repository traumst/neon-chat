/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./static/html/**/**"],
  theme: {
    extend: {
      colors: {
        dark: {
          DEFAULT: '#000000',
          100: '#111111',
          200: '#222222',
        }, // /dark
        light: {
          DEFAULT: '#FFFFFF',
          100: '#EEEEEE',
          200: '#DDDDDD',
        }, // /light
        crazy: {
          DEFAULT: '#6b46c1',
          100: '#452959',
          200: '#a45fba',
        }, // /crazy
      }, // /colors
    }, // /extend
  }, // /theme
  plugins: [],
}

