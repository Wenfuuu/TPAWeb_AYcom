<script lang="ts">
  import { onMount } from "svelte";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import { fetchUser, type User } from "../lib/model/user";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import {
    fetchAllNotifications,
    fetchMentionNotifications,
    type Notif,
  } from "../lib/model/notification";
  import NotificationContent from "../lib/components/NotificationContent.svelte";
  import { toastStore } from "../lib/stores/toast";
  import axios from "axios";

  let isDarkMode = true;
  let user: User | null = null;
  let notifications: Notif[] = [];

  let socket: WebSocket;

  function setupWebSocket(userId: string) {
    if (socket) {
      socket.close();
    }

    socket = new WebSocket(`ws://localhost:50056/ws?user_id=${userId}`);

    socket.onmessage = async (event) => {
      const wsMessage = JSON.parse(event.data);
      // console.log(wsMessage.type);
      switch (wsMessage.type) {
        case "notification": {
          const notification = wsMessage.payload;
          // console.log("Received notification:", notification);
          try {
            const allUsers = await fetchUser();
            const sender = allUsers.find((u) => u.id === notification.sender_id);
            const profile_picture = sender ? sender.profile_picture : "";

            // Add profile picture to notification
            const notificationWithPicture = {
              ...notification,
              profile_picture
            };

            toastStore.info(notification.content);
            notifications = [notificationWithPicture, ...notifications];
          } catch (error) {
            if (axios.isAxiosError(error)) {
              const message = error.response?.data?.error || error.message;
              toastStore.error(message);
            } else if (error instanceof Error) {
              toastStore.error(error.message);
            } else {
              toastStore.error("An unexpected error occurred!");
            }
          }
          break;
        }

        default:
          // console.log("Unknown message type:", wsMessage.type);
      }
    };

    socket.onopen = () => {
      // console.log("WebSocket connection established");
    };

    socket.onerror = (err) => {
      if (err instanceof Error) {
        toastStore.error(err.message);
      } else {
        toastStore.error("An unexpected error occurred!");
      }
    };

    // socket.onclose = () => {
    //   console.log("WebSocket connection closed. Reconnecting...");
    //   setTimeout(() => setupWebSocket(userId), 5000); // Reconnect after 5 seconds
    // };
  }

  async function refreshData(
    page: number = 1,
    currTab: "all" | "mentions" = "all"
  ) {
    if (!user) return;

    if (currTab === "all") notifications = await fetchAllNotifications(user.id);
    else if (currTab === "mentions")
      notifications = await fetchMentionNotifications(user.id);

    // console.log(notifications);
  }

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  onMount(async () => {
    await checkAuth(navigate);

    const usr = await getAuthUser();
    if (usr) {
      user = usr;
      setupWebSocket(user.id);
    }

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      const prefersDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute(
        "data-theme",
        prefersDark ? "dark" : "light"
      );
    }

    await refreshData();
  });
</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />
  <NotificationContent {notifications} {refreshData} {user} />
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
