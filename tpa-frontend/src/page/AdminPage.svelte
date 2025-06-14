<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import { checkAdmin, getAuthUser } from "../lib/utils/auth";
  import { toastStore } from "../lib/stores/toast";
  import { acceptPremium, acceptReport, fetchPremiumRequests, fetchReports, fetchUser, fetchUserById, rejectPremium, rejectReport, setUserBanStatus, type PremiumRequest, type Reports, type User } from "../lib/model/user";
  import { acceptCreateCommunity, createCommunityCategory, deleteCommCategory, fetchCommunityCategory, fetchCreateCommunity, rejectCreateCommunity, updateCommCategory, type CommunityCategory, type CreateCommunityRequest } from "../lib/model/community";
  import axios from "axios";
  import { createThreadCategory, deleteThrdCategory, fetchThreadCategory, updateThrdCategory, type ThreadCategory } from "../lib/model/thread";
  import { sendMessageNotif } from "../lib/model/websocket";

  // Tab management
  const TABS = {
    USERS: "users",
    NEWSLETTER: "newsletter",
    COMMUNITY_REQUESTS: "community_requests",
    PREMIUM_REQUESTS: "premium_requests",
    REPORTS: "reports",
    THREAD_CATEGORIES: "thread_categories",
    COMMUNITY_CATEGORIES: "community_categories"
  };

  let user: User | null = null;
  let currentTab = TABS.USERS;
  let isDarkMode = true;
  let isLoading = false;
  let socket: WebSocket;

  // Data states
  let users: User[] = [];
  let communityRequests: CreateCommunityRequest[] = [];
  let premiumRequests: PremiumRequest[] = [];
  let reportRequests: Reports[] = [];
  let threadCategories: ThreadCategory[] = [];
  let communityCategories: CommunityCategory[] = [];

  // Form states
//   let newsletterSubject = "";
  let newsletterContent = "";
  let newThreadCategory = "";
  let newCommunityCategory = "";
  let editingThreadCategory: { id: string; name: string } | null = null;
  let editingCommunityCategory: { id: string; name: string } | null = null;

  // Modal states
  let showModal = false;
  let modalTitle = "";
  let modalContent = "";
  let modalAction: (() => void) | null = null;
  // const modalType = "confirm";

  function setupWebSocket(userId: string) {
    if (socket) {
      socket.close();
    }

    socket = new WebSocket(`ws://localhost:50056/ws?user_id=${userId}`);

    socket.onmessage = (event) => {
      const wsMessage = JSON.parse(event.data);
      // console.log(wsMessage.type);
      switch (wsMessage.type) {
        case "notification": {
          const notification = wsMessage.payload;
          toastStore.info(notification.content);
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

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  async function loadUsers() {
    isLoading = true;
    try {
      users = await fetchUser();
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  async function loadCommunityRequests() {
    isLoading = true;
    try {
      communityRequests = await fetchCreateCommunity();
    } catch (error) {
        if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
        } else if (error instanceof Error) {
            toastStore.error(error.message);
        }
    } finally {
      isLoading = false;
    }
  }

  async function loadPremiumRequests() {
    isLoading = true;
    try {
      premiumRequests = await fetchPremiumRequests();
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  async function loadReportRequests() {
    isLoading = true;
    try {
      reportRequests = await fetchReports();
      // In a real app, this would be an API call
      // await new Promise(resolve => setTimeout(resolve, 800));
      // reportRequests = [
      //   { id: "1", reporter: "user1", reported_user: "user2", reason: "Inappropriate content", status: "pending", created_at: "2023-05-13T16:20:00Z" },
      //   { id: "2", reporter: "user3", reported_user: "user5", reason: "Harassment", status: "pending", created_at: "2023-05-14T08:30:00Z" }
      // ];
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  async function loadThreadCategories() {
    isLoading = true;
    try {
      threadCategories = await fetchThreadCategory();
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  async function loadCommunityCategories() {
    isLoading = true;
    try {
      communityCategories = await fetchCommunityCategory();
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  function changeTab(tab: string) {
    currentTab = tab;

    // Load data based on selected tab
    switch (tab) {
      case TABS.USERS:
        loadUsers();
        break;
      case TABS.COMMUNITY_REQUESTS:
        loadCommunityRequests();
        break;
      case TABS.PREMIUM_REQUESTS:
        loadPremiumRequests();
        break;
      case TABS.REPORTS:
        loadReportRequests();
        break;
      case TABS.THREAD_CATEGORIES:
        loadThreadCategories();
        break;
      case TABS.COMMUNITY_CATEGORIES:
        loadCommunityCategories();
        break;
    }
  }

  async function toggleUserBan(userId: string, currentBanStatus: boolean) {
    showConfirmModal(
      currentBanStatus ? "Unban User" : "Ban User",
      `Are you sure you want to ${currentBanStatus ? "unban" : "ban"} this user?`,
      async () => {
        try {
          isLoading = true;
          await setUserBanStatus(userId, !currentBanStatus);

          // Update local state
          users = await fetchUser();

          toastStore.success(`User ${currentBanStatus ? "unbanned" : "banned"} successfully`);
        } catch (error) {
            if (axios.isAxiosError(error)) {
                const message = error.response?.data?.error || error.message;
                toastStore.error(message);
            } else if (error instanceof Error) {
                toastStore.error(error.message);
            }
        } finally {
          isLoading = false;
          closeModal();
        }
      }
    );
  }

  async function sendNewsletter() {
    if(!user) return;
    if (!newsletterContent.trim()) {
      toastStore.error("Please fill in content");
      return;
    }

    try {
      isLoading = true;

      const subscribedUsers = users.filter(user => user.is_subscribed);
      for (const subsUser of subscribedUsers) {
        const email = subsUser.email;

        const notif = {
            recipient_id: subsUser.id,
            sender_id: user.id,
            content: `Newsletter: ${newsletterContent}`,
            category: "message",
            email: email
        };

        const notifMessage = {
            type: "notification",
            payload: notif
        };

        socket.send(JSON.stringify(notifMessage));
      }

      toastStore.success("Newsletter sent successfully");
    //   newsletterSubject = "";
      newsletterContent = "";
    } catch (error) {
        if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
        } else if (error instanceof Error) {
            toastStore.error(error.message);
        }
    } finally {
      isLoading = false;
    }
  }

  async function handleCommunityRequest(request: CreateCommunityRequest, action: "accept" | "reject") {
    showConfirmModal(
      `${action === "accept" ? "Accept" : "Reject"} Community Request`,
      `Are you sure you want to ${action} this community request?`,
      async () => {
        try {
          isLoading = true;

          //update
          if(action === "accept") {
            await acceptCreateCommunity(request.id);
            if(!user) return;
            const recipientUser = await fetchUserById(request.user_id);
            sendMessageNotif(socket, request.user_id, user.id, "Your create community request has been accepteed.", recipientUser.email);
          }else {
            await rejectCreateCommunity(request.id);
          }

          loadCommunityRequests();

          toastStore.success(`Community request ${action}ed successfully`);
        } catch (error) {
            if (axios.isAxiosError(error)) {
                const message = error.response?.data?.error || error.message;
                toastStore.error(message);
            } else if (error instanceof Error) {
                toastStore.error(error.message);
            }
        } finally {
          isLoading = false;
          closeModal();
        }
      }
    );
  }

  async function handlePremiumRequest(request: PremiumRequest, action: "accept" | "reject") {
    showConfirmModal(
      `${action === "accept" ? "Accept" : "Reject"} Premium Request`,
      `Are you sure you want to ${action} this premium request?`,
      async () => {
        try {
          isLoading = true;

          //update
          if(action === "accept") {
            await acceptPremium(request.id);
            if(!user) return;
            const recipientUser = await fetchUserById(request.user_id);
            sendMessageNotif(socket, request.user_id, user.id, "Your premium request has been accepted.", recipientUser.email);
          }else {
            await rejectPremium(request.id);
          }

          loadPremiumRequests();

          toastStore.success(`Premium request ${action}ed successfully`);
        } catch (error) {
            if (axios.isAxiosError(error)) {
                const message = error.response?.data?.error || error.message;
                toastStore.error(message);
            } else if (error instanceof Error) {
                toastStore.error(error.message);
            }
        } finally {
          isLoading = false;
          closeModal();
        }
      }
    );
  }

  async function handleReport(report_id: string, action: "accept" | "reject") {
    showConfirmModal(
      `${action === "accept" ? "Accept" : "Reject"} Report Request`,
      `Are you sure you want to ${action} this report request?`,
      async () => {
        try {
          isLoading = true;

          //update
          if(action === "accept") {
            await acceptReport(report_id);
          }else {
            await rejectReport(report_id);
          }

          loadReportRequests();

          toastStore.success(`Report request ${action}ed successfully`);

          // In a real app, this would trigger an email notification
          if (action === "accept") {
            // console.log("Email notification would be sent for premium request acceptance");

          }
        } catch (error) {
            if (axios.isAxiosError(error)) {
                const message = error.response?.data?.error || error.message;
                toastStore.error(message);
            } else if (error instanceof Error) {
                toastStore.error(error.message);
            }
        } finally {
          isLoading = false;
          closeModal();
        }
      }
    );
  }

  async function addThreadCategory() {
    if (!newThreadCategory.trim()) {
      toastStore.error("Category name cannot be empty");
      return;
    }

    try {
      isLoading = true;

      await createThreadCategory(newThreadCategory);
      loadThreadCategories();

      newThreadCategory = "";
      toastStore.success("Thread category added successfully");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  async function addCommunityCategory() {
    if (!newCommunityCategory.trim()) {
      toastStore.error("Category name cannot be empty");
      return;
    }

    try {
      isLoading = true;

      await createCommunityCategory(newCommunityCategory);

      loadCommunityCategories();

      newCommunityCategory = "";
      toastStore.success("Community category added successfully");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  function editThreadCategory(category: { id: string; name: string }) {
    editingThreadCategory = { ...category };
  }

  function editCommunityCategory(category: { id: string; name: string }) {
    editingCommunityCategory = { ...category };
  }

  async function updateThreadCategory() {
    if (!editingThreadCategory) return;
    if (!editingThreadCategory.name.trim()) {
      toastStore.error("Category name cannot be empty");
      return;
    }

    try {
      isLoading = true;

      await updateThrdCategory(editingThreadCategory.id, editingThreadCategory.name);
      loadThreadCategories();

      editingThreadCategory = null;
      toastStore.success("Thread category updated successfully");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  async function updateCommunityCategory() {
    if (!editingCommunityCategory) return;
    if (!editingCommunityCategory.name.trim()) {
      toastStore.error("Category name cannot be empty");
      return;
    }

    try {
      isLoading = true;

      await updateCommCategory(editingCommunityCategory.id, editingCommunityCategory.name);
      loadCommunityCategories();

      editingCommunityCategory = null;
      toastStore.success("Community category updated successfully");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    } finally {
      isLoading = false;
    }
  }

  function deleteThreadCategory(id: string) {
    showConfirmModal(
      "Delete Thread Category",
      "Are you sure you want to delete this thread category? This action cannot be undone.",
      async () => {
        try {
          isLoading = true;

          await deleteThrdCategory(id);
          loadThreadCategories();

          toastStore.success("Thread category deleted successfully");
        } catch (error) {
          if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
          } else if (error instanceof Error) {
            toastStore.error(error.message);
          }
        } finally {
          isLoading = false;
          closeModal();
        }
      }
    );
  }

  function deleteCommunityCategory(id: string) {
    showConfirmModal(
      "Delete Community Category",
      "Are you sure you want to delete this community category? This action cannot be undone.",
      async () => {
        try {
          isLoading = true;

          await deleteCommCategory(id);
          loadCommunityCategories();

          toastStore.success("Community category deleted successfully");
        } catch (error) {
          if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
          } else if (error instanceof Error) {
            toastStore.error(error.message);
          }
        } finally {
          isLoading = false;
          closeModal();
        }
      }
    );
  }

  function showConfirmModal(title: string, content: string, action: () => void) {
    modalTitle = title;
    modalContent = content;
    modalAction = action;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    modalAction = null;
  }

  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit"
    }).format(date);
  }

  onMount(async () => {
    await checkAdmin(navigate);

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
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute("data-theme", prefersDark ? "dark" : "light");
    }

    // Load initial data
    loadUsers();
  });
</script>

<main class="admin-page">
  <div class="top-controls">
    <h1>Admin Dashboard</h1>
    <button
      type="button"
      class="theme-toggle"
      aria-label="Toggle theme"
      on:click={toggleTheme}
    >
      {#if isDarkMode}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="5"></circle>
          <line x1="12" y1="1" x2="12" y2="3"></line>
          <line x1="12" y1="21" x2="12" y2="23"></line>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
          <line x1="1" y1="12" x2="3" y2="12"></line>
          <line x1="21" y1="12" x2="23" y2="12"></line>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
        </svg>
      {:else}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
        </svg>
      {/if}
    </button>
  </div>

  <div class="tabs">
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.USERS}
      on:click={() => changeTab(TABS.USERS)}
    >
      Users
    </button>
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.NEWSLETTER}
      on:click={() => changeTab(TABS.NEWSLETTER)}
    >
      Newsletter
    </button>
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.COMMUNITY_REQUESTS}
      on:click={() => changeTab(TABS.COMMUNITY_REQUESTS)}
    >
      Community Requests
    </button>
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.PREMIUM_REQUESTS}
      on:click={() => changeTab(TABS.PREMIUM_REQUESTS)}
    >
      Premium Requests
    </button>
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.REPORTS}
      on:click={() => changeTab(TABS.REPORTS)}
    >
      Reports
    </button>
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.THREAD_CATEGORIES}
      on:click={() => changeTab(TABS.THREAD_CATEGORIES)}
    >
      Thread Categories
    </button>
    <button type="button"
      class="tab-button"
      class:active={currentTab === TABS.COMMUNITY_CATEGORIES}
      on:click={() => changeTab(TABS.COMMUNITY_CATEGORIES)}
    >
      Community Categories
    </button>
  </div>

  <div class="content">
    {#if isLoading}
      <div class="loading-container">
        <div class="spinner"></div>
        <p>Loading...</p>
      </div>
    {:else}
      <!-- Users Tab -->
      {#if currentTab === TABS.USERS}
        <div class="section">
          <h2>User Management</h2>

          {#if users.length === 0}
            <p class="empty-state">No users found</p>
          {:else}
            <div class="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Username</th>
                    <th>Email</th>
                    <th>Status</th>
                    <th>Premium</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {#each users as user (user.id)}
                    <tr>
                      <td>{user.username}</td>
                      <td>{user.email}</td>
                      <td>
                        <span class="status-badge {user.is_banned ? "banned" : "active"}">
                          {user.is_banned ? "Banned" : "Active"}
                        </span>
                      </td>
                      <td>
                        <span class="status-badge {user.is_premium ? "premium" : "standard"}">
                          {user.is_premium ? "Premium" : "Standard"}
                        </span>
                      </td>
                      <td>
                        <button type="button"
                          class="action-button {user.is_banned ? "unban" : "ban"}"
                          on:click={() => toggleUserBan(user.id, user.is_banned)}
                        >
                          {user.is_banned ? "Unban" : "Ban"}
                        </button>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>

      <!-- Newsletter Tab -->
      {:else if currentTab === TABS.NEWSLETTER}
        <div class="section">
          <h2>Send Newsletter</h2>
          <p class="section-description">
            Send a newsletter to all subscribed users.
          </p>

          <form on:submit|preventDefault={sendNewsletter}>
            <div class="form-group">
              <label for="newsletter-content">Content</label>
              <textarea
                id="newsletter-content"
                bind:value={newsletterContent}
                placeholder="Write your newsletter content here..."
                rows="10"
                required
              ></textarea>
            </div>

            <button
              type="submit"
              class="primary-button"
              disabled={isLoading}
            >
              {#if isLoading}
                <span class="spinner-small"></span>
                Sending...
              {:else}
                Send Newsletter
              {/if}
            </button>
          </form>
        </div>

      <!-- Community Requests Tab -->
      {:else if currentTab === TABS.COMMUNITY_REQUESTS}
        <div class="section">
          <h2>Community Creation Requests</h2>

          {#if communityRequests.length === 0}
            <p class="empty-state">No pending community requests</p>
          {:else}
            <div class="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Creator</th>
                    <th>Description</th>
                    <th>Date</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {#each communityRequests as request (request.id)}
                    <tr>
                      <td>{request.name}</td>
                      <td>{request.creator_name}</td>
                      <td class="description-cell">{request.description}</td>
                      <td>{formatDate(request.created_at)}</td>
                      <td>
                          <div class="action-buttons">
                            <button type="button"
                              class="action-button accept"
                              on:click={() => handleCommunityRequest(request, "accept")}
                            >
                              Accept
                            </button>
                            <button type="button"
                              class="action-button reject"
                              on:click={() => handleCommunityRequest(request, "reject")}
                            >
                              Reject
                            </button>
                          </div>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>

      <!-- Premium Requests Tab -->
      {:else if currentTab === TABS.PREMIUM_REQUESTS}
        <div class="section">
          <h2>Premium Verification Requests</h2>

          {#if premiumRequests.length === 0}
            <p class="empty-state">No pending premium requests</p>
          {:else}
            <div class="table-container">
              <table>
                <thead>
                  <tr>
                    <th>User</th>
                    <th>User ID</th>
                    <th>Reason</th>
                    <th>Date</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {#each premiumRequests as request (request.id)}
                    <tr>
                      <td>
                        <img class="face-pic" src={request.face_picture} alt="">
                      </td>
                      <td>{request.user_id}</td>
                      <td class="description-cell">{request.reason}</td>
                      <td>{formatDate(request.created_at)}</td>
                      <td>
                        <div class="action-buttons">
                            <button type="button"
                                class="action-button accept"
                                on:click={() => handlePremiumRequest(request, "accept")}
                            >
                                Accept
                            </button>
                            <button type="button"
                                class="action-button reject"
                                on:click={() => handlePremiumRequest(request, "reject")}
                            >
                                Reject
                            </button>
                        </div>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>

      <!-- Reports Tab -->
      {:else if currentTab === TABS.REPORTS}
        <div class="section">
          <h2>User Reports</h2>

          {#if reportRequests.length === 0}
            <p class="empty-state">No pending reports</p>
          {:else}
            <div class="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Reporter ID</th>
                    <th>Reported ID</th>
                    <th>Reason</th>
                    <th>Date</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {#each reportRequests as report (report.id)}
                    <tr>
                      <td>{report.user_id}</td>
                      <td>{report.reported_id}</td>
                      <td class="description-cell">{report.reason}</td>
                      <td>{formatDate(report.created_at)}</td>
                      <td>
                          <div class="action-buttons">
                            <button type="button"
                              class="action-button accept"
                              on:click={() => handleReport(report.id, "accept")}
                              title="Accept report and ban user"
                            >
                              Accept
                            </button>
                            <button type="button"
                              class="action-button reject"
                              on:click={() => handleReport(report.id, "reject")}
                            >
                              Reject
                            </button>
                          </div>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>

      <!-- Thread Categories Tab -->
      {:else if currentTab === TABS.THREAD_CATEGORIES}
        <div class="section">
          <h2>Thread Categories</h2>

          <div class="category-form">
            <div class="form-group">
              <label for="new-thread-category">Add New Category</label>
              <div class="input-with-button">
                <input
                  type="text"
                  id="new-thread-category"
                  bind:value={newThreadCategory}
                  placeholder="Enter thread category name"
                />
                <button
                  type="button"
                  class="primary-button"
                  on:click={addThreadCategory}
                  disabled={isLoading || !newThreadCategory.trim()}
                >
                  Add
                </button>
              </div>
            </div>
          </div>

          {#if threadCategories.length === 0}
            <p class="empty-state">No thread categories found</p>
          {:else}
            <div class="categories-list">
              {#each threadCategories as category (category.id)}
                <div class="category-item">
                  {#if editingThreadCategory && editingThreadCategory.id === category.id}
                    <div class="category-edit">
                      <input
                        type="text"
                        bind:value={editingThreadCategory.name}
                        placeholder="Category name"
                      />
                      <div class="edit-actions">
                        <button type="button"
                          class="action-button save"
                          on:click={updateThreadCategory}
                          disabled={isLoading}
                        >
                          Save
                        </button>
                        <button type="button"
                          class="action-button cancel"
                          on:click={() => editingThreadCategory = null}
                          disabled={isLoading}
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  {:else}
                    <span class="category-name">{category.name}</span>
                    <div class="category-actions">
                      <button type="button"
                        class="action-button edit"
                        on:click={() => editThreadCategory(category)}
                      >
                        Edit
                      </button>
                      <button type="button"
                        class="action-button delete"
                        on:click={() => deleteThreadCategory(category.id)}
                      >
                        Delete
                      </button>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}
        </div>

      <!-- Community Categories Tab -->
      {:else if currentTab === TABS.COMMUNITY_CATEGORIES}
        <div class="section">
          <h2>Community Categories</h2>

          <div class="category-form">
            <div class="form-group">
              <label for="new-community-category">Add New Category</label>
              <div class="input-with-button">
                <input
                  type="text"
                  id="new-community-category"
                  bind:value={newCommunityCategory}
                  placeholder="Enter community category name"
                />
                <button
                  type="button"
                  class="primary-button"
                  on:click={addCommunityCategory}
                  disabled={isLoading || !newCommunityCategory.trim()}
                >
                  Add
                </button>
              </div>
            </div>
          </div>

          {#if communityCategories.length === 0}
            <p class="empty-state">No community categories found</p>
          {:else}
            <div class="categories-list">
              {#each communityCategories as category (category.id)}
                <div class="category-item">
                  {#if editingCommunityCategory && editingCommunityCategory.id === category.id}
                    <div class="category-edit">
                      <input
                        type="text"
                        bind:value={editingCommunityCategory.name}
                        placeholder="Category name"
                      />
                      <div class="edit-actions">
                        <button type="button"
                          class="action-button save"
                          on:click={updateCommunityCategory}
                          disabled={isLoading}
                        >
                          Save
                        </button>
                        <button type="button"
                          class="action-button cancel"
                          on:click={() => editingCommunityCategory = null}
                          disabled={isLoading}
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  {:else}
                    <span class="category-name">{category.name}</span>
                    <div class="category-actions">
                      <button type="button"
                        class="action-button edit"
                        on:click={() => editCommunityCategory(category)}
                      >
                        Edit
                      </button>
                      <button type="button"
                        class="action-button delete"
                        on:click={() => deleteCommunityCategory(category.id)}
                      >
                        Delete
                      </button>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    {/if}
  </div>

  <!-- Confirmation Modal -->
  {#if showModal}
  <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
    <div class="modal-overlay" on:click={closeModal}>
      <div class="modal" on:click|stopPropagation>
        <div class="modal-header">
          <h3>{modalTitle}</h3>
          <button type="button" class="close-button" aria-label="close button" on:click={closeModal}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-content">
          <p>{modalContent}</p>
        </div>

        <div class="modal-actions">
          <button class="secondary-button" type="button" on:click={closeModal}>Cancel</button>
          <button type="button"
            class="primary-button"
            on:click={modalAction}
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="spinner-small"></span>
              Processing...
            {:else}
              Confirm
            {/if}
          </button>
        </div>
      </div>
    </div>
  {/if}
</main>

<style>
  .admin-page {
    width: 100%;
    min-height: 100vh;
    background-color: var(--background);
    color: var(--text-primary);
    padding-bottom: 40px;
  }

  .top-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid var(--border-color);
    background-color: var(--card-background);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  h1 {
    font-size: 24px;
    font-weight: 700;
    margin: 0;
  }

  .theme-toggle {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .theme-toggle:hover {
    background-color: var(--hover-background);
  }

  .tabs {
    display: flex;
    overflow-x: auto;
    background-color: var(--card-background);
    border-bottom: 1px solid var(--border-color);
    position: sticky;
    top: 65px;
    z-index: 9;
  }

  .tab-button {
    padding: 16px 20px;
    background: none;
    border: none;
    color: var(--text-tertiary);
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    position: relative;
    transition: color 0.2s;
  }

  .tab-button.active {
    color: var(--button-primary);
    font-weight: 700;
  }

  .tab-button.active::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background-color: var(--button-primary);
    border-radius: 4px 4px 0 0;
  }

  .content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .section {
    background-color: var(--card-background);
    border-radius: 16px;
    padding: 24px;
    border: 1px solid var(--border-color);
    margin-bottom: 24px;
  }

  h2 {
    font-size: 20px;
    font-weight: 700;
    margin: 0 0 16px 0;
  }

  .section-description {
    color: var(--text-secondary);
    margin-bottom: 24px;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(var(--button-primary-rgb, 29, 155, 240), 0.3);
    border-radius: 50%;
    border-top-color: var(--button-primary);
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 16px;
  }

  .spinner-small {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 1s ease-in-out infinite;
    margin-right: 8px;
  }

  .face-pic {
    max-width: 100px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .empty-state {
    text-align: center;
    color: var(--text-tertiary);
    padding: 40px;
  }

  .table-container {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    text-align: left;
    padding: 12px 16px;
    font-weight: 600;
    color: var(--text-secondary);
    border-bottom: 1px solid var(--border-color);
  }

  td {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
  }

  tr:last-child td {
    border-bottom: none;
  }

  .description-cell {
    max-width: 300px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .status-badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 9999px;
    font-size: 12px;
    font-weight: 600;
  }

  .status-badge.active {
    background-color: rgba(52, 199, 89, 0.1);
    color: #34c759;
  }

  .status-badge.banned {
    background-color: rgba(255, 69, 58, 0.1);
    color: #ff453a;
  }

  .status-badge.premium {
    background-color: rgba(29, 155, 240, 0.1);
    color: #1d9bf0;
  }

  .status-badge.standard {
    background-color: rgba(142, 142, 147, 0.1);
    color: #8e8e93;
  }

  .action-button {
    padding: 6px 12px;
    border-radius: 9999px;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    border: none;
    transition: background-color 0.2s;
  }

  .action-button.ban {
    background-color: rgba(255, 69, 58, 0.1);
    color: #ff453a;
  }

  .action-button.ban:hover {
    background-color: rgba(255, 69, 58, 0.2);
  }

  .action-button.unban {
    background-color: rgba(52, 199, 89, 0.1);
    color: #34c759;
  }

  .action-button.unban:hover {
    background-color: rgba(52, 199, 89, 0.2);
  }

  .action-button.accept {
    background-color: rgba(52, 199, 89, 0.1);
    color: #34c759;
  }

  .action-button.accept:hover {
    background-color: rgba(52, 199, 89, 0.2);
  }

  .action-button.reject {
    background-color: rgba(255, 69, 58, 0.1);
    color: #ff453a;
  }

  .action-button.reject:hover {
    background-color: rgba(255, 69, 58, 0.2);
  }

  .action-button.edit {
    background-color: rgba(29, 155, 240, 0.1);
    color: #1d9bf0;
  }

  .action-button.edit:hover {
    background-color: rgba(29, 155, 240, 0.2);
  }

  .action-button.delete {
    background-color: rgba(255, 69, 58, 0.1);
    color: #ff453a;
  }

  .action-button.delete:hover {
    background-color: rgba(255, 69, 58, 0.2);
  }

  .action-button.save {
    background-color: rgba(52, 199, 89, 0.1);
    color: #34c759;
  }

  .action-button.save:hover {
    background-color: rgba(52, 199, 89, 0.2);
  }

  .action-button.cancel {
    background-color: rgba(142, 142, 147, 0.1);
    color: var(--text-primary);
  }

  .action-button.cancel:hover {
    background-color: rgba(142, 142, 147, 0.2);
  }

  .action-buttons {
    display: flex;
    gap: 8px;
  }

  .form-group {
    margin-bottom: 20px;
  }

  label {
    display: block;
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  input, textarea {
    width: 100%;
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 12px;
    color: var(--text-primary);
    font-size: 15px;
    transition: border-color 0.2s;
    box-sizing: border-box;
    font-family: inherit;
  }

  textarea {
    resize: vertical;
    min-height: 100px;
  }

  input:focus, textarea:focus {
    border-color: var(--button-primary);
    outline: none;
  }

  .input-with-button {
    display: flex;
    gap: 8px;
  }

  .input-with-button input {
    flex: 1;
  }

  .primary-button {
    background-color: var(--button-primary);
    border: none;
    color: white;
    padding: 10px 20px;
    border-radius: 9999px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .primary-button:hover:not(:disabled) {
    background-color: var(--button-primary-hover);
  }

  .primary-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .secondary-button {
    background: none;
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    padding: 10px 20px;
    border-radius: 9999px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .secondary-button:hover {
    background-color: var(--hover-background);
  }

  .categories-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 24px;
  }

  .category-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background-color: var(--hover-background);
    border-radius: 8px;
    border: 1px solid var(--border-color);
  }

  .category-name {
    font-weight: 500;
  }

  .category-actions {
    display: flex;
    gap: 8px;
  }

  .category-edit {
    display: flex;
    width: 100%;
    gap: 12px;
    align-items: center;
  }

  .category-edit input {
    flex: 1;
  }

  .edit-actions {
    display: flex;
    gap: 8px;
  }

  .category-form {
    background-color: var(--hover-background);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 24px;
    border: 1px solid var(--border-color);
  }

  .category-form .form-group:last-child {
    margin-bottom: 0;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
  }

  .modal {
    background-color: var(--card-background);
    border-radius: 16px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    overflow: hidden;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-color);
  }

  .modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 700;
  }

  .close-button {
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .close-button:hover {
    background-color: var(--hover-background);
  }

  .modal-content {
    padding: 20px;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 20px;
    border-top: 1px solid var(--border-color);
  }

  @media (max-width: 768px) {
    .tabs {
      padding: 0 8px;
    }

    .tab-button {
      padding: 12px 10px;
      font-size: 13px;
    }

    .section {
      padding: 16px;
    }

    th, td {
      padding: 8px;
    }

    .description-cell {
      max-width: 150px;
    }

    .action-buttons {
      flex-direction: column;
      gap: 4px;
    }

    .category-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .category-actions {
      align-self: flex-end;
    }

    .category-edit {
      flex-direction: column;
      align-items: stretch;
    }
  }
</style>