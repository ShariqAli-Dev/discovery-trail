/** @type {import('tailwindcss').Config} */
module.exports = {
  // content: ["./html/**/*.{html,js,templ}"],
  content: ["../html/**/*.{html,js,templ}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["autumn", "nord", "dracula"],
  },
};
