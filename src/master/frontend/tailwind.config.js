/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        transparent: 'transparent',
        current: 'currentColor',
        'mainBlue': '#003049',
        'mainRed': '#D62828',
        'mainOrange': '#F77F00',
        'mainGreen': '#70C1B3',
        'mainDark': '#2E0014'
      },
    },
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
