/** @type {import('tailwindcss').Config} */
module.exports = {
  important: true,
  content: ["./assets/**/*.{css,js}", "./templates/**/*.html"],
  theme: { extend: {} },
  plugins: [require("@tailwindcss/forms")],
};
