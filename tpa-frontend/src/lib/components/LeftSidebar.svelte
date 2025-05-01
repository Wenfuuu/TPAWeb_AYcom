<script lang="ts">
  import { onMount } from "svelte";
  import { type User } from "../model/user";
  import { logout } from "../utils/auth";
  import { navigate } from "svelte-routing";
  import ComposePost from "./ComposePost.svelte";
  import Modal from "./Modal.svelte";

  export let toggleTheme: () => void;
  export let isDarkMode: boolean;
  export let user: User | null;
  export let currentActiveTab: "for-you" | "following";

  export let refreshData: (page: number) => Promise<void>;
  export let refreshFollowingData: (page: number) => Promise<void>;

  let showComposeModal = false;

  let showProfileDropdown = false;
  let dropdownRef: HTMLDivElement;

  const navItems = [
    { name: "Home", icon: "home" },
    { name: "Explore", icon: "search" },
    { name: "Notifications", icon: "bell" },
    { name: "Messages", icon: "mail" },
    { name: "Bookmarks", icon: "bookmark" },
    { name: "Communities", icon: "users" },
    { name: "Premium", icon: "star" },
    { name: "Profile", icon: "user" },
    { name: "More", icon: "more-horizontal" }
  ];

  // Close dropdown when clicking outside
  function handleClickOutside(event: MouseEvent) {
    if (dropdownRef && !dropdownRef.contains(event.target as Node)) {
      showProfileDropdown = false;
    }
  }

  function toggleProfileDropdown() {
    showProfileDropdown = !showProfileDropdown;
  }

  function handleLogout() {
    logout(navigate);
  }

  function navigateHome() {
    navigate("/home");
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter" || event.key === " ") {
      toggleProfileDropdown();
      event.preventDefault();
    }
  }

  function openComposeModal() {
    showComposeModal = true;
  }

  function closeComposeModal() {
    showComposeModal = false;
  }

  // Function called when ComposePost inside the modal finishes posting
  async function handleModalPostCreated() {
    console.log("Post created via modal. Closing modal.");
    if(currentActiveTab === "for-you") await refreshData(1);
    else refreshFollowingData(1);

    closeComposeModal();
  }

  onMount(() => {
    document.addEventListener("click", handleClickOutside);
    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  });
</script>

<div class="sidebar left-sidebar">

  <div class="sidebar-content">
    <button type="button" class="logo" on:click={navigateHome}>
      {#if isDarkMode}
        <img src="/src/assets/AY-dark.jpg" alt="Logo" class="logo-img">
      {:else}
        <img src="/src/assets/AY-light.jpg" alt="Logo" class="logo-img">
      {/if}
    </button>

    <nav class="navigation">
      {#each navItems as item (item.name)}
        <a href="#{item.name.toLowerCase()}" class="nav-item">
          <span class="icon">
            {#if item.icon === "home"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
            {:else if item.icon === "search"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><path d="m21 21-4.3-4.3"></path></svg>
            {:else if item.icon === "bell"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"></path><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"></path></svg>
            {:else if item.icon === "mail"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="16" x="2" y="4" rx="2"></rect><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"></path></svg>
            {:else if item.icon === "bookmark"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m19 21-7-4-7 4V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"></path></svg>
            {:else if item.icon === "users"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M22 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            {:else if item.icon === "star"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            {:else if item.icon === "user"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            {:else if item.icon === "more-horizontal"}
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"></circle><circle cx="5" cy="12" r="1"></circle></svg>
            {/if}
          </span>
          <span class="nav-text">{item.name}</span>
        </a>
      {/each}
    </nav>

    <button type="button" class="post-button" on:click={openComposeModal}>Post</button>

    {#if user}
      <div class="user-profile-container" bind:this={dropdownRef}>
        <div
          class="user-profile"
          role="button"
          tabindex="0"
          on:click={toggleProfileDropdown}
          on:keydown={handleKeyDown}
          aria-expanded={showProfileDropdown}
          aria-haspopup="true"
        >
          <img src={user.profile_picture || "/placeholder.svg"} alt={user.name || "User"} class="avatar" />
          <div class="user-info">
            <div class="user-name">{user.name || "User"}</div>
            <div class="user-username">{user.username || "@user"}</div>
          </div>

          <div class="theme-toggle">
            <button type="button" on:click|stopPropagation={toggleTheme} aria-label="Toggle theme">
              {#if isDarkMode}
                <!-- Sun icon for light mode -->
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
                <!-- Moon icon for dark mode -->
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
                </svg>
              {/if}
            </button>
          </div>

          <button type="button" class="more-options" aria-label="More options" on:click|stopPropagation={toggleProfileDropdown}>
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="1"></circle>
              <circle cx="19" cy="12" r="1"></circle>
              <circle cx="5" cy="12" r="1"></circle>
            </svg>
          </button>
        </div>

        {#if showProfileDropdown}
          <div class="profile-dropdown">
            <div class="dropdown-options">
              <button type="button" class="dropdown-option" on:click={handleLogout}>
                Log out {user.username || "@user"}
              </button>
            </div>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<Modal show={showComposeModal} on:close={closeComposeModal}>
  {#if user} <!-- Only render ComposePost if user is available -->
    <ComposePost
      {user}
      refreshData={handleModalPostCreated}
      {currentActiveTab} />
  {/if}
</Modal>

<style>
  .left-sidebar {
    width: 25%;
    border-right: 1px solid var(--border-color);
    position: sticky;
    top: 0;
    height: 100vh;
    overflow-y: auto;
  }

  .sidebar-content {
    display: flex;
    justify-content: right;
    flex-direction: column;
    height: 100%;
    padding: 0 12px;
  }

  .logo {
    background-color: var(--background);
    border: none;
    padding: 12px;
    margin-top: 4px;
    margin-bottom: 8px;
    width: fit-content;
    border-radius: 50%;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .logo:hover {
    background-color: var(--hover-background);
  }

  .logo-img {
    max-width: 50px;
  }

  .navigation {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    padding: 12px;
    border-radius: 9999px;
    color: var(--text-primary);
    text-decoration: none;
    font-size: 18px;
    font-weight: 400;
    transition: background-color 0.2s;
  }

  .nav-item:hover {
    background-color: var(--hover-background);
  }

  .icon {
    margin-right: 16px;
    display: flex;
    align-items: center;
  }

  .post-button {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 16px 0;
    font-size: 17px;
    font-weight: 700;
    cursor: pointer;
    margin-top: 16px;
    margin-bottom: 24px;
    transition: background-color 0.2s;
  }

  .post-button:hover {
    background-color: var(--button-primary-hover);
  }

  .theme-toggle {
    display: flex;
    justify-content: center;
    margin: 16px 0;
  }

  .theme-toggle button {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    padding: 8px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .theme-toggle button:hover {
    background-color: var(--hover-background);
  }

  .user-profile-container {
    position: relative;
    margin-bottom: 16px;
  }

  /* Updated styles for user-profile as a button */
  .user-profile {
    position: relative;
    display: flex;
    align-items: center;
    padding: 12px 0;
    border-radius: 9999px;
    cursor: pointer;
    transition: background-color 0.2s;
    background: none;
    border: none;
    width: 100%;
    text-align: left;
    color: var(--text-primary);
  }

  .user-profile:hover {
    background-color: var(--hover-background);
  }

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    margin-right: 12px;
  }

  .user-info {
    flex: 1;
    overflow: hidden;
  }

  .user-name {
    font-weight: 700;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-username {
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .more-options {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    padding: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .profile-dropdown {
    position: absolute;
    bottom: calc(100% + 10px);
    left: 0;
    width: 80%;
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .dropdown-options {
    padding: 12px 0;
  }

  .dropdown-option {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 15px;
    font-weight: 500;
    padding: 12px 16px;
    text-align: left;
    cursor: pointer;
    transition: background-color 0.2s;
    width: 100%;
  }

  .dropdown-option:hover {
    background-color: var(--hover-background);
  }

  @media (max-width: 768px) {
    .left-sidebar {
      width: 100%;
      height: auto;
      position: fixed;
      bottom: 0;
      top: auto;
      border-right: none;
      border-top: 1px solid var(--border-color);
      z-index: 100;
      background-color: var(--background);
    }

    .sidebar-content {
      flex-direction: row;
      justify-content: space-around;
      padding: 8px 0;
    }

    .navigation {
      flex-direction: row;
      margin-bottom: 0;
      flex: 1;
      justify-content: space-around;
    }

    .logo, .post-button, .theme-toggle, .user-profile, .nav-text {
      display: none;
    }

    .nav-item {
      padding: 8px;
    }

    .profile-dropdown {
      bottom: auto;
      top: 100%;
      right: 0;
      left: auto;
    }
  }
</style>