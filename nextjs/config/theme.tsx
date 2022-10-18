import { createTheme, responsiveFontSizes } from "@mui/material/styles";

import { red } from "@mui/material/colors";

function getTheme(prefersDarkMode: boolean) {
  return responsiveFontSizes(
    createTheme({
      palette: {
        mode: prefersDarkMode ? "dark" : "light",
      },
      components: {
        MuiUseMediaQuery: {
          defaultProps: {
            noSsr: true,
          },
        },
      },
    })
  );
}

export default getTheme;
