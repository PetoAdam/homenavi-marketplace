/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#0b1117",
        slate: "#1a2430",
        panel: "#1f2a36",
        accent: "#39aa79",
        mist: "#dbe2ea"
      },
      fontFamily: {
        display: ["Space Grotesk", "Instrument Sans", "sans-serif"],
        body: ["Instrument Sans", "Space Grotesk", "sans-serif"]
      },
      boxShadow: {
        soft: "0 18px 40px rgba(10, 16, 24, 0.35)"
      }
    }
  },
  plugins: []
};
