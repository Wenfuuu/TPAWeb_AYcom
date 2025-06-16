<script lang="ts">
  import { onMount } from "svelte";
  import type { User } from "../model/user";

  export let user: User | null = null;

  let activeTab = "security";

  let isPrivateAccount = false;
  let fontSize = 16;
  let fontColor = "#000000";

  let notificationSettings = {
    like: true,
    repost: true,
    follow: true,
    mentions: true,
    community: true,
  };

  let blockedAccounts = [
    { id: "1", username: "user1", displayName: "John Doe", avatar: "" },
    { id: "2", username: "user2", displayName: "Jane Smith", avatar: "" },
  ];

  onMount(() => {
    // Load saved settings from localStorage
    const savedFontSize = localStorage.getItem("fontSize");
    if (savedFontSize) fontSize = parseInt(savedFontSize);

    const savedFontColor = localStorage.getItem("fontColor");
    if (savedFontColor) fontColor = savedFontColor;

    const savedNotifications = localStorage.getItem("notificationSettings");
    if (savedNotifications)
      notificationSettings = JSON.parse(savedNotifications);
  });

  function switchTab(tab: string) {
    activeTab = tab;
  }

  function saveDisplaySettings() {
    localStorage.setItem("fontSize", fontSize.toString());
    localStorage.setItem("fontColor", fontColor);
  }

  function toggleAccountPrivacy() {
    isPrivateAccount = !isPrivateAccount;
    // TODO: Call API to update account privacy
  }

  function deactivateAccount() {
    if (
      confirm(
        "Are you sure you want to deactivate your account? This action cannot be undone."
      )
    ) {
      // TODO: Call API to deactivate account
      alert("Account deactivation functionality not implemented yet.");
    }
  }

  function unblockUser(userId: string) {
    blockedAccounts = blockedAccounts.filter(
      (account) => account.id !== userId
    );
    // TODO: Call API to unblock user
  }

  function saveNotificationSettings() {
    localStorage.setItem(
      "notificationSettings",
      JSON.stringify(notificationSettings)
    );
    // TODO: Call API to save notification preferences
  }

  const tabs = [
    { id: "security", name: "Security", icon: "🛡️" },
    { id: "display", name: "Display", icon: "🖥️" },
    { id: "account", name: "Your Account", icon: "👤" },
    { id: "blocked", name: "Blocked Accounts", icon: "🚫" },
    { id: "notifications", name: "Notification Preferences", icon: "🔔" },
  ];
  $: userDisplayName = user?.username || "User";

  $: headerTitle = userDisplayName
    ? `Settings for @${userDisplayName}`
    : "Settings";
</script>

<main class="settings-content">
  <div class="settings-header">
    <h1>{headerTitle}</h1>
  </div>

  <div class="settings-container">
    <!-- Tab Navigation -->
    <div class="tab-navigation">
      {#each tabs as tab (tab.id)}
        <button
          type="button"
          class="tab-button"
          class:active={activeTab === tab.id}
          on:click={() => switchTab(tab.id)}
        >
          <span class="tab-icon">{tab.icon}</span>
          {tab.name}
        </button>
      {/each}
    </div>

    <!-- Tab Content -->
    <div class="tab-content">
      {#if activeTab === "security"}
        <div class="setting-section">
          <h2>Account Privacy</h2>
          <div class="setting-item">
            <div class="setting-info">
              <h3>Private Account</h3>
              <p>
                When your account is private, only people you approve can see
                your posts and follow you.
              </p>
            </div>
            <label class="toggle-switch">
              <input
                type="checkbox"
                bind:checked={isPrivateAccount}
                on:change={toggleAccountPrivacy}
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      {:else if activeTab === "display"}
        <div class="setting-section">
          <h2>Customize Your Experience</h2>

          <div class="setting-item">
            <div class="setting-info">
              <h3>Font Size</h3>
              <p>Adjust the font size for better readability.</p>
            </div>
            <div class="font-size-controls">
              <input
                type="range"
                min="12"
                max="24"
                bind:value={fontSize}
                on:input={saveDisplaySettings}
              />
              <span class="font-size-value">{fontSize}px</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <h3>Font Color</h3>
              <p>Choose your preferred text color.</p>
            </div>
            <input
              type="color"
              bind:value={fontColor}
              on:change={saveDisplaySettings}
              class="color-picker"
            />
          </div>

          <div class="preview-text">
            This is a preview of your text settings.
          </div>
        </div>
      {:else if activeTab === "account"}
        <div class="setting-section">
          <h2>Account Management</h2>

          <div class="setting-item danger-zone">
            <div class="setting-info">
              <h3>Deactivate Account</h3>
              <p>
                Temporarily deactivate your account. You can reactivate it
                anytime by logging back in.
              </p>
            </div>
            <button
              type="button"
              class="danger-button"
              on:click={deactivateAccount}
            >
              Deactivate Account
            </button>
          </div>
        </div>
      {:else if activeTab === "blocked"}
        <div class="setting-section">
          <h2>Blocked Accounts</h2>

          {#if blockedAccounts.length === 0}
            <div class="empty-state">
              <p>You haven't blocked any accounts yet.</p>
            </div>
          {:else}
            <div class="blocked-accounts-list">
              {#each blockedAccounts as account (account.id)}
                <div class="blocked-account-item">
                  <div class="account-info">
                    <div class="avatar">
                      {account.displayName.charAt(0)}
                    </div>
                    <div class="account-details">
                      <div class="display-name">{account.displayName}</div>
                      <div class="username">@{account.username}</div>
                    </div>
                  </div>
                  <button
                    type="button"
                    class="unblock-button"
                    on:click={() => unblockUser(account.id)}
                  >
                    Unblock
                  </button>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {:else if activeTab === "notifications"}
        <div class="setting-section">
          <h2>Notification Preferences</h2>
          <p>Choose which notifications you'd like to receive.</p>

          <div class="notification-settings">
            <div class="notification-item">
              <div class="notification-info">
                <h3>Likes</h3>
                <p>When someone likes your posts</p>
              </div>
              <label class="toggle-switch">
                <input
                  type="checkbox"
                  bind:checked={notificationSettings.like}
                  on:change={saveNotificationSettings}
                />
                <span class="slider"></span>
              </label>
            </div>

            <div class="notification-item">
              <div class="notification-info">
                <h3>Reposts</h3>
                <p>When someone reposts your content</p>
              </div>
              <label class="toggle-switch">
                <input
                  type="checkbox"
                  bind:checked={notificationSettings.repost}
                  on:change={saveNotificationSettings}
                />
                <span class="slider"></span>
              </label>
            </div>

            <div class="notification-item">
              <div class="notification-info">
                <h3>Follows</h3>
                <p>When someone follows you</p>
              </div>
              <label class="toggle-switch">
                <input
                  type="checkbox"
                  bind:checked={notificationSettings.follow}
                  on:change={saveNotificationSettings}
                />
                <span class="slider"></span>
              </label>
            </div>

            <div class="notification-item">
              <div class="notification-info">
                <h3>Mentions</h3>
                <p>When someone mentions you in a post</p>
              </div>
              <label class="toggle-switch">
                <input
                  type="checkbox"
                  bind:checked={notificationSettings.mentions}
                  on:change={saveNotificationSettings}
                />
                <span class="slider"></span>
              </label>
            </div>

            <div class="notification-item">
              <div class="notification-info">
                <h3>Community</h3>
                <p>Community updates and announcements</p>
              </div>
              <label class="toggle-switch">
                <input
                  type="checkbox"
                  bind:checked={notificationSettings.community}
                  on:change={saveNotificationSettings}
                />
                <span class="slider"></span>
              </label>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</main>

<style>
  .settings-content {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    margin: 0 auto;
    overflow-y: auto;
  }

  .settings-header {
    position: sticky;
    top: 0;
    background-color: rgba(var(--background), 0.95);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--border-color);
    padding: 1.25rem 1.5rem;
    z-index: 100;
  }

  .settings-header h1 {
    font-size: 1.5rem;
    font-weight: 800;
    margin: 0;
    color: var(--text-primary);
    letter-spacing: -0.025em;
  }

  .settings-container {
    display: flex;
    flex-direction: column;
  }

  .tab-navigation {
    display: flex;
    flex-direction: column;
    border-bottom: 1px solid var(--border-color);
    background-color: var(--card-background);
  }

  .tab-button {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem 1.5rem;
    background: none;
    border: none;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    text-align: left;
    font-size: 1rem;
    font-weight: 500;
    position: relative;
  }

  .tab-button:hover {
    background-color: var(--hover-background);
    color: var(--text-primary);
  }

  .tab-button.active {
    color: var(--button-primary);
    background-color: var(--hover-background);
    font-weight: 600;
  }

  .tab-button.active::after {
    content: "";
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background-color: var(--button-primary);
    border-radius: 2px 0 0 2px;
  }

  .tab-icon {
    font-size: 1.25rem;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .tab-content {
    background-color: var(--card-background);
  }

  .setting-section {
    padding: 2rem 1.5rem;
  }

  .setting-section h2 {
    font-size: 1.75rem;
    font-weight: 700;
    margin: 0 0 0.75rem 0;
    color: var(--text-primary);
    letter-spacing: -0.025em;
  }

  .setting-section > p {
    color: var(--text-secondary);
    margin-bottom: 2rem;
    font-size: 1rem;
    line-height: 1.5;
  }

  .setting-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 1.75rem 0;
    border-bottom: 1px solid var(--border-color);
    gap: 1.5rem;
    transition: background-color 0.2s ease;
  }

  .setting-item:last-child {
    border-bottom: none;
  }

  .setting-info {
    flex: 1;
  }

  .setting-info h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin: 0 0 0.5rem 0;
    color: var(--text-primary);
  }

  .setting-info p {
    font-size: 0.9rem;
    color: var(--text-secondary);
    margin: 0;
    line-height: 1.5;
  }

  .toggle-switch {
    position: relative;
    display: inline-block;
    width: 52px;
    height: 28px;
    flex-shrink: 0;
  }

  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--border-color);
    transition: 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 28px;
    box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 22px;
    width: 22px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 50%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  input:checked + .slider {
    background-color: var(--button-primary);
  }

  input:checked + .slider:before {
    transform: translateX(24px);
  }

  .slider:hover {
    box-shadow:
      inset 0 1px 3px rgba(0, 0, 0, 0.15),
      0 0 0 3px rgba(29, 155, 240, 0.1);
  }

  .font-size-controls {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-shrink: 0;
  }
  .font-size-controls input[type="range"] {
    width: 140px;
    height: 6px;
    border-radius: 3px;
    background: var(--border-color);
    outline: none;
    -webkit-appearance: none;
    appearance: none;
  }
  .font-size-controls input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: var(--button-primary);
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .font-size-controls input[type="range"]::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: var(--button-primary);
    cursor: pointer;
    border: none;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .font-size-value {
    font-weight: 600;
    color: var(--text-primary);
    min-width: 45px;
    font-size: 0.9rem;
    background-color: var(--hover-background);
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
  }

  .color-picker {
    width: 52px;
    height: 32px;
    border: 2px solid var(--border-color);
    border-radius: 8px;
    cursor: pointer;
    flex-shrink: 0;
    transition: border-color 0.2s ease;
  }

  .color-picker:hover {
    border-color: var(--button-primary);
  }

  .preview-text {
    margin-top: 1.5rem;
    padding: 1.25rem;
    border: 1px solid var(--border-color);
    border-radius: 12px;
    background-color: var(--hover-background);
    font-style: italic;
    color: var(--text-secondary);
  }

  .danger-zone {
    border: 2px solid #ef4444;
    border-radius: 12px;
    padding: 1.75rem;
    background: linear-gradient(
      135deg,
      rgba(239, 68, 68, 0.05) 0%,
      rgba(239, 68, 68, 0.02) 100%
    );
    position: relative;
    overflow: hidden;
  }

  .danger-zone::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, #ef4444, #dc2626);
  }

  .danger-button {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    color: white;
    border: none;
    padding: 0.875rem 1.75rem;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.95rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
  }

  .danger-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
  }

  .danger-button:active {
    transform: translateY(0);
  }

  .empty-state {
    text-align: center;
    padding: 4rem 1rem;
    color: var(--text-secondary);
  }

  .empty-state p {
    font-size: 1.1rem;
    margin: 0;
  }

  .blocked-accounts-list {
    display: flex;
    flex-direction: column;
    gap: 0;
    background-color: var(--card-background);
    border-radius: 12px;
    border: 1px solid var(--border-color);
    overflow: hidden;
  }

  .blocked-account-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
    transition: background-color 0.2s ease;
  }

  .blocked-account-item:last-child {
    border-bottom: none;
  }

  .blocked-account-item:hover {
    background-color: var(--hover-background);
  }

  .account-info {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .avatar {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--button-primary) 0%, #1a8cd8 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: 600;
    font-size: 1.1rem;
    flex-shrink: 0;
  }

  .account-details {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .display-name {
    font-weight: 600;
    color: var(--text-primary);
    font-size: 1rem;
  }

  .username {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .unblock-button {
    background-color: var(--button-primary);
    color: white;
    border: none;
    padding: 0.625rem 1.25rem;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.875rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    flex-shrink: 0;
  }

  .unblock-button:hover {
    background-color: var(--button-primary-hover);
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(29, 155, 240, 0.3);
  }

  .notification-settings {
    display: flex;
    flex-direction: column;
    background-color: var(--card-background);
    border-radius: 12px;
    border: 1px solid var(--border-color);
    overflow: hidden;
  }

  .notification-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--border-color);
    transition: background-color 0.2s ease;
  }

  .notification-item:last-child {
    border-bottom: none;
  }

  .notification-item:hover {
    background-color: var(--hover-background);
  }

  .notification-info {
    flex: 1;
  }

  .notification-info h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin: 0 0 0.375rem 0;
    color: var(--text-primary);
  }

  .notification-info p {
    font-size: 0.9rem;
    color: var(--text-secondary);
    margin: 0;
    line-height: 1.4;
  }

  /* Responsive design */
  @media (max-width: 768px) {
    .settings-content {
      max-width: 100%;
      border-left: none;
      border-right: none;
    }

    .setting-section {
      padding: 1.5rem 1rem;
    }

    .setting-item {
      flex-direction: column;
      align-items: stretch;
      gap: 1rem;
    }

    .font-size-controls,
    .toggle-switch,
    .color-picker {
      align-self: flex-start;
    }

    .blocked-account-item,
    .notification-item {
      padding: 1rem;
    }
  }

  .tab-content {
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .tab-button:focus-visible {
    outline: 2px solid var(--button-primary);
    outline-offset: 2px;
  }

  .toggle-switch input:focus-visible + .slider {
    box-shadow: 0 0 0 3px rgba(29, 155, 240, 0.3);
  }

  .danger-button:focus-visible,
  .unblock-button:focus-visible {
    outline: 2px solid var(--button-primary);
    outline-offset: 2px;
  }

  .color-picker:focus-visible {
    outline: 2px solid var(--button-primary);
    outline-offset: 2px;
  }

  .font-size-controls input[type="range"]:focus {
    outline: none;
  }

  .font-size-controls input[type="range"]:focus::-webkit-slider-thumb {
    box-shadow:
      0 2px 4px rgba(0, 0, 0, 0.2),
      0 0 0 3px rgba(29, 155, 240, 0.3);
  }

  .setting-item {
    transition: background-color 0.2s ease;
  }

  .setting-item:hover {
    background-color: rgba(var(--hover-background), 0.5);
    border-radius: 8px;
    margin: 0 -0.5rem;
    padding-left: 2rem;
    padding-right: 2rem;
  }
</style>
