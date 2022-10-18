import type { NextPage } from "next";
import Head from "next/head";

import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import TextField from "@mui/material/TextField";
import { Button } from "@mui/material";
import Typography from "@mui/material/Typography";
import { useTheme } from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";

const baseBoxStyle = {
  display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
};

const mainContainerStyle = {
  textAlign: "center",
};

const introTextContainerStyle = {
  textAlign: "left",
  mr: "0em",
  ml: "5%",
  mt: "2em",
  px: "0",
};

const introTextStyle = {
  maxWidth: "20rem",
};

const textBoxStyle = {
  display: "inline-flex",
  alignItems: "center",
  justifyContent: "center",
  my: "2em",
  flexDirection: "row",
  width: "100%",
  // width
};

const createPageButtonStyle = {
  zIndex: 1,
  display: "absolute",
  left: "-10em",
  top: "1em",
};

const Home: NextPage = ({}) => {
  const theme = useTheme();
  const sm = useMediaQuery(theme.breakpoints.up("sm"));

  return (
    <>
      <Container sx={introTextContainerStyle}>
        <Typography variant="h4" sx={introTextStyle}>
          Share your favorite resources, all in one place
        </Typography>
      </Container>
      <Container sx={mainContainerStyle}>
        <Head>
          <title>LinkShare</title>
        </Head>

        <Box sx={textBoxStyle}>
          <TextField
            label="linkshare.dev/"
            helperText="Enter an optional custom URL"
            sx={{ width: "60%", margin: "5em 0em 2em 0" }}
          />
          <Button variant={"contained"} sx={createPageButtonStyle}>
            Create Page
          </Button>
        </Box>
      </Container>
    </>
  );
};

export default Home;
