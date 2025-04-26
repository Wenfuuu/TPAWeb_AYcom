<script lang="ts">
  import { Router, Link, Route } from "svelte-routing";
  import LandingPage from "./page/LandingPage.svelte";
  import RegisterPage from "./page/RegisterPage.svelte";

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

</Router>

<style>
  
</style>
