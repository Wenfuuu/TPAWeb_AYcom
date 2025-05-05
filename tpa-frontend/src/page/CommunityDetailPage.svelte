<script lang="ts">
  import { onMount } from "svelte";
  import type { User, UserThread } from "../lib/model/user";
  import { fetchCommunityDetailWithCreator, fetchModerator, fetchTopMember, type Community, type CommunityMember } from "../lib/model/community";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import { fetchCommunityThreadWithUser, fetchTopCommunityThreadWithUser, type Thread } from "../lib/model/thread";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import CommunityDetailContent from "../lib/components/CommunityDetailContent.svelte";

  let isDarkMode = true; // Default to dark mode
  let user: User | null = null;
  let community: Community;
  let threads: (Thread & { userData: UserThread })[] = [];
  let topMembers: CommunityMember[] = [];
  let moderators: CommunityMember[] = [];

    async function refreshData(page: number,
        currTab: "top" | "latest" | "media" | "about" | "manage" = "latest"
    ){
        if(user) {
            let data;
            console.log(currTab);
            if(currTab === "latest") data = await fetchCommunityThreadWithUser(page, user.id, community.id);
            else if(currTab === "top") {
                console.log("fetching top");
                data = await fetchTopCommunityThreadWithUser(page, user.id, community.id);
            }
            else data = await fetchCommunityThreadWithUser(page, user.id, community.id);

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

        if(user) community = await fetchCommunityDetailWithCreator(user.id, communityId);

        topMembers = await fetchTopMember(community.id);
        // console.log(topMembers);
        moderators = await fetchModerator(community.id);
        // console.log(moderators);

        refreshData(1);
    });

</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />
  <CommunityDetailContent {community} {threads} {user} {refreshData}
  {topMembers} {moderators} />
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