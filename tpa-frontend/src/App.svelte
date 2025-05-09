<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import LandingPage from "./page/LandingPage.svelte";
  import RegisterPage from "./page/RegisterPage.svelte";
  import HomePage from "./page/HomePage.svelte";
  import LoginPage from "./page/LoginPage.svelte";
  import ThreadDetailPage from "./page/ThreadDetailPage.svelte";
  import BookmarkPage from "./page/BookmarkPage.svelte";
  import PremiumPage from "./page/PremiumPage.svelte";
  import CommunityPage from "./page/CommunityPage.svelte";
  import CommunityDetailPage from "./page/CommunityDetailPage.svelte";
  import ProfilePage from "./page/ProfilePage.svelte";
  import CreateCommunityPage from "./page/CreateCommunityPage.svelte";
  import ToastContainer from "./lib/components/ToastContainer.svelte";
  import MessagePage from "./page/MessagePage.svelte";
  import NotificationPage from "./page/NotificationPage.svelte";
  import AdminPage from "./page/AdminPage.svelte";
  import ForgottenPage from "./page/ForgottenPage.svelte";
  import ExplorePage from "./page/ExplorePage.svelte";
  import SettingsPage from "./page/SettingsPage.svelte";

  export let url = "";
</script>

<ToastContainer />
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

  <Route path="/bookmarks">
    <BookmarkPage />
  </Route>

  <Route path="/premium">
    <PremiumPage />
  </Route>

  <Route path="/communities">
    <CommunityPage />
  </Route>

  <Route path="/community-detail/:id">
    <CommunityDetailPage />
  </Route>

   <Route path="/create-community">
    <CreateCommunityPage />
  </Route>

  <Route path="/profile/:id">
    <ProfilePage />
  </Route>

  <Route path="/messages">
    <MessagePage />
  </Route>

  <Route path="/notifications">
    <NotificationPage />
  </Route>

  <Route path="/explore">
    <ExplorePage />
  </Route>

  <Route path="/admin">
    <AdminPage />
  </Route>

  <Route path="/forgotten">
    <ForgottenPage />
  </Route>

  <Route path="/settings">
    <SettingsPage />
  </Route>

</Router>

<style>

</style>
