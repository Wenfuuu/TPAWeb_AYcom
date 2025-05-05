<script lang="ts">
  import { onMount } from "svelte";
  import type { User } from "../lib/model/user";
  import { fetchCommunityDetail, type Community } from "../lib/model/community";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import { fetchCommunityThreadWithUser } from "../lib/model/thread";

  let isDarkMode = true; // Default to dark mode
  let user: User | null = null;
  let community: Community;
  let threads: (Thread & { userData: UserThread })[] = [];

    async function refreshData(page: number,
        currTab: "top" | "latest" | "media" | "about" | "manage" = "top"
    ){
        if(user) {
            const data = await fetchCommunityThreadWithUser(page, user.id, community.id);
            if(data) threads = data;
            console.log(threads);
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

        // get current community
        const currentUrl = window.location.href;
        const urlParts = currentUrl.split("/");
        const communityId = urlParts[urlParts.length - 1];
        community = await fetchCommunityDetail(communityId);
        console.log(community);

        // refresh threads based on page
        refreshData(1);
    });

</script>



<style>

</style>