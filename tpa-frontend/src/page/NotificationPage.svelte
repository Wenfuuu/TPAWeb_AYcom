<script lang="ts">
    import { onMount } from "svelte";
    import LeftSidebar from "../lib/components/LeftSidebar.svelte";
    import RightSidebar from "../lib/components/RightSidebar.svelte";
    import type { User } from "../lib/model/user";
    import { checkAuth, getAuthUser } from "../lib/utils/auth";
    import { navigate } from "svelte-routing";
    import { fetchNotifications, type Notification } from "../lib/model/notification";
    import NotificationContent from "../lib/components/NotificationContent.svelte";

    let isDarkMode = true;
    let user: User | null = null;
    let notifications: Notification[] = [];

    async function refreshData(page: number = 1,
        currTab: "all" | "mentions" = "all"
    ) {
        if(!user) return;
        notifications = await fetchNotifications(user.id);

        console.log(notifications);
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

        await refreshData();
    });

</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />
  <!-- <NotificationContent {refreshData} /> -->
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