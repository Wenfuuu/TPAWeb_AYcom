<script lang="ts">
  import { onMount } from "svelte";
  import { fetchDetailedLikedThread, fetchDetailedUserReply, fetchDetailedUserThread, fetchThreadMedia, type ThreadMedia } from "../lib/model/thread";
  import { fetchProfileUser, type User } from "../lib/model/user";
  import type { Thread } from "../lib/model/thread";
  import type { UserThread } from "../lib/model/user";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import ProfileContent from "../lib/components/ProfileContent.svelte";

    let isDarkMode = true;
    let user: User | null = null;
    let threads: (Thread & { userData: UserThread })[] = [];
    let threadMedia: ThreadMedia[] = [];
    let isUser = false;
    let userProfile: User | null = null;

    async function refreshData(page: number,
        currTab: "posts" | "replies" | "likes" | "media" = "posts"
    ){
        if(user) {
            const currentUrl = window.location.href;
            const urlParts = currentUrl.split("/");
            const profileId = urlParts[urlParts.length - 1];

            userProfile = await fetchProfileUser(profileId);
            if(user.id === userProfile.id) isUser = true;
            else isUser = false;

            let data;
            if(currTab === "posts") data = await fetchDetailedUserThread(userProfile.id);
            else if(currTab === "replies") data = await fetchDetailedUserReply(userProfile.id);
            else if(currTab === "likes") data = await fetchDetailedLikedThread(userProfile.id);
            else data = await fetchDetailedUserThread(userProfile.id);

            if(data) threads = data;

            threadMedia = await fetchThreadMedia(userProfile.id);
        }
    }

    function toggleTheme() {
      isDarkMode = !isDarkMode;
      const newTheme = isDarkMode ? "dark" : "light";
      localStorage.setItem("theme", newTheme);
      document.documentElement.setAttribute("data-theme", newTheme);
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

        refreshData(1);
    });

</script>

<div class="home-container">
    <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />
    <ProfileContent {user} {threads} {refreshData}
    {threadMedia} {isUser} {userProfile} />
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