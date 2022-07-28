/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/pages/**/*{js,ts,jsx,tsx}",
    "./src/components/**/*{js,ts,jsx,tsx}",
  ],
  theme: {
    screens: {
      xs: "375px",
      sm: "600px",
      md: "900px",
      lg: "1200px",
      xl: "1536px",
    },
    fontFamily: {
      sans: ["Arial", "sans-serif"],
      serif: ["Garamond", "serif"],
    },
  },
  plugins: [],
};
