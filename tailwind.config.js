/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./assets/**/*.{css,js}", "./templates/**/*.html"],
  theme: { extend: {} },
  plugins: [require("@tailwindcss/forms")],
};
