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
  let currId = "";
  const idStack: string[] = [];

  async function refreshData(page: number, fromBack = false): Promise<void> {
    console.log("stackk:  " + idStack);

    const currentUrl = window.location.href;
    const urlParts = currentUrl.split("/");
    const threadId = urlParts[urlParts.length - 1];

    // Check if we've navigated to a different thread
    if (currId !== "" && currId !== threadId && !fromBack) {
        // Store previous thread before updating current
        // prevId = currId;
        idStack.push(currId);
        currId = threadId;

        // Load the new thread data without page reload
        await loadThreadData(threadId, page);
        return;
    } else if (currId === "") {
        // First load of any thread
        currId = threadId;
    }

    // Load thread data
    await loadThreadData(threadId, page);
  }

    async function loadThreadData(threadId: string, page: number): Promise<void> {
        if (user) {
            // Fetch thread details
            const data = await fetchThreadDetailWithUser(threadId, user.id);
            if (data) thread = data;

            // Fetch thread replies
            const replyData = await fetchThreadRepliesWithUser(page, threadId, user.id);
            if (replyData) replies = replyData;
            else replies = [];

            console.log("Thread loaded:", thread);
            console.log("Replies loaded:", replies.length);
        }
    }

    function handleBackNavigation() {
        if (idStack.length > 0) {
            const prevId = idStack.pop();
            if (prevId) {
                navigate(`/thread-detail/${prevId}`);
                refreshData(1, true);
            }
        } else {
            navigate("/home");
        }

        // console.log("aaaaaaa" + prevId);
        // if (prevId) {
        //     navigate(`/thread-detail/${prevId}`);
        //     refreshData(1);
        // } else {
        //     navigate("/home");
        // }
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
   <ThreadDetailContent {thread} {replies} {refreshData}
    {loadThreadData} {user} {handleBackNavigation} />
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