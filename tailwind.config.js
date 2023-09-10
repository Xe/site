/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./tmpl/*.html"], // This is where your HTML templates / JSX files are located
  theme: {
    extend: {
      fontFamily: {
        sans: ["Iosevka Aile Iaso", "sans-serif"],
        mono: ["Iosevka Curly Iaso", "monospace"],
        serif: ["Iosevka Etoile Iaso", "serif"],
      },
      colors: {
        bg: {
          hard: "#f9f5d7",
          soft: "#f2e5bc",
          0: "#fbf1c7",
          1: "#ebdbb2",
          2: "#d5c4a1",
          3: "#bdae93",
          4: "#a89984",
        },
        bgDark: {
          hard: "#1d2021",
          soft: "#32302f",
          0: "#282828",
          1: "#3c3836",
          2: "#504945",
          3: "#665c54",
          4: "#7c6f64",
        },
        fg: {
          0: "#282828",
          1: "#3c3836",
          2: "#504945",
          3: "#665c54",
          4: "#7c6f64",
        },
        fgDark: {
          0: "#fbf1c7",
          1: "#ebdbb2",
          2: "#d5c4a1",
          3: "#bdae93",
          4: "#a89984",
        },
        red: {
          dark: "#9d0006",
          light: "#cc241d",
        },
        redDark: {
          dark: "#cc241d",
          light: "#fb4934",
        },
        green: {
          light: "#98971a",
          dark: "#79740e",
        },
        greenDark: {
          dark: "#98971a",
          light: "#b8bb26",
        },
        yellow: {
          light: "#d79921",
          dark: "#b57614",
        },
        yellowDark: {
          light: "#d79921",
          dark: "#fabd2f",
        },
        blue: {
          light: "#458588",
          dark: "#076678",
        },
        blueDark: {
          dark: "#458588",
          light: "#83a598",
        },
        purple: {
          light: "#b16286",
          dark: "#8f3f71",
        },
        purpleDark: {
          dark: "#b16286",
          light: "#d3869b",
        },
        aqua: {
          light: "#689d6a",
          dark: "#427b58",
        },
        aquaDark: {
          dark: "#689d6a",
          light: "#8ec07c",
        },
        orange: {
          light: "#d65d0e",
          dark: "#af3a03",
        },
        orangeDark: {
          dark: "#d65d0e",
          light: "#fe8019",
        },
        gray: {
          light: "#928374",
          dark: "#7c6f64",
        },
        grayDark: {
          dark: "#928374",
          light: "#a89984",
        },
      }
    },
  },
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};