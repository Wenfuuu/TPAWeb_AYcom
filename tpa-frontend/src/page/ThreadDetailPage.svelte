<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import type { User, UserThread } from "../lib/model/user";
  import { fetchThreadDetailWithUser, fetchThreadRepliesWithUser, type Thread } from "../lib/model/thread";
  import ThreadDetailContent from "../lib/components/ThreadDetailContent.svelte";

  let isDarkMode = true; // Default to dark mode
  let user: User | null = null;
  let thread: (Thread & { userData: UserThread });
  let replies: (Thread & { userData: UserThread })[] = [];
  const currentActiveTab: "for-you" | "following" = "for-you";
  let prevId = "";

  async function refreshData(page: number): Promise<void> {
    const currentUrl = window.location.href;
    const urlParts = currentUrl.split("/");
    const threadId = urlParts[urlParts.length - 1];
    if(prevId != "" && prevId != threadId) location.reload();
    prevId = threadId;

    if(user){
        const data = await fetchThreadDetailWithUser(threadId, user.id);
        if(data) thread = data;

        const replyData = await fetchThreadRepliesWithUser(page, threadId, user.id);
        if(replyData) replies = replyData;
    }
    console.log(thread);
    console.log(replies);
    console.log(replies.length);
  }

  onMount(async() => {
    await checkAuth(navigate);

    const usr = await getAuthUser();
    if(usr) user = usr;

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute("data-theme", prefersDark ? "dark" : "light");
    }

    await refreshData(1);
  });

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }
</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} refreshFollowingData={refreshData}
  {currentActiveTab} />
  <!-- <HomeContent {threads} {user} {refreshData} {refreshFollowingData} {currentActiveTab} /> -->
   <ThreadDetailContent {thread} {replies} {refreshData} {user} />
  <RightSidebar />
</div>

<style>
  .home-container {
    display: flex;
    width: 100%;
    height: 100vh;
    background-color: var(--background);
    color: var(--text-primary);
  }
</style>