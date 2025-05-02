<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import LandingPage from "./page/LandingPage.svelte";
  import RegisterPage from "./page/RegisterPage.svelte";
  import HomePage from "./page/HomePage.svelte";
  import LoginPage from "./page/LoginPage.svelte";
  import ThreadDetailPage from "./page/ThreadDetailPage.svelte";

  // You need to export the URL prop to make router work
  export let url = "";
</script>

<Router {url}>
  <Route path="/">
    <LandingPage />
  </Route>

  <Route path="/register">
    <RegisterPage />
  </Route>

  <Route path="/home">
    <HomePage />
  </Route>

  <Route path="/login">
    <LoginPage />
  </Route>

  <Route path="/thread-detail/:id">
    <ThreadDetailPage />
  </Route>

</Router>

<style>

</style>
