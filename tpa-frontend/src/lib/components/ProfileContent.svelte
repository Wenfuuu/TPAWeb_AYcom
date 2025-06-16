<script lang="ts">
  import { onMount } from "svelte";
  import type { Thread, ThreadMedia } from "../model/thread";
  import { followUser, unfollowUser, type UserThread } from "../model/user";
  import type { User } from "../model/user";
  import { navigate } from "svelte-routing";
  import ThreadFeed from "./ThreadFeed.svelte";
  import { toastStore } from "../stores/toast";

  export let user: User | null;
  export let userProfile: User | null;
  export let threads: (Thread & { userData: UserThread })[] = [];
  export let refreshData: (
    page: number,
    currTab: "posts" | "replies" | "likes" | "media"
  ) => Promise<void>;
  export let threadMedia: ThreadMedia[] = [];
  export let isUser = false;

  let activeTab: "posts" | "replies" | "likes" | "media" = "posts";
  let page = 1;
  let showProfileModal = false;

  function openProfileModal() {
    showProfileModal = true;
  }

  function closeProfileModal() {
    showProfileModal = false;
  }

  function handleModalClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      closeProfileModal();
    }
  }
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      closeProfileModal();
    }
  }

  onMount(() => {
    refreshData(page, activeTab);

    // Add global keyboard listener for modal
    const handleGlobalKeydown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && showProfileModal) {
        closeProfileModal();
      }
    };

    document.addEventListener("keydown", handleGlobalKeydown);

    return () => {
      document.removeEventListener("keydown", handleGlobalKeydown);
    };
  });
</script>

{#if userProfile}
  <main class="content-area">
    <div class="profile-header">
      <div class="banner-container">
        <img
          src={userProfile?.banner}
          alt="{userProfile?.name} banner"
          class="profile-banner"
        />
        <button
          type="button"
          class="back-button"
          aria-label="Back to communities"
          on:click={() => navigate("/communities")}
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M19 12H5"
              stroke="white"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M12 19L5 12L12 5"
              stroke="white"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>
      </div>
      <div class="profile-info">
        <button
          type="button"
          class="profile-logo-button"
          on:click={openProfileModal}
          aria-label="View profile picture"
        >
          <img
            src={userProfile?.profile_picture}
            alt="{userProfile?.name} profile picture"
            class="profile-logo"
          />
        </button>
        <div class="profile-details">
          <h1 class="profile-name">{userProfile?.name}</h1>
          <span class="profile-username">@{userProfile?.username}</span>
          <p class="profile-bio">{userProfile?.bio}</p>
          <p class="profile-joined">
            Joined: {new Date(userProfile.created_at).toLocaleDateString()}
          </p>
        </div>
        {#if !isUser}
          <button
            type="button"
            class="follow-button"
            on:click={async () => {
              if (!user) return;
              if (!userProfile.is_followed) {
                await followUser(user.id, userProfile.id);
                toastStore.success("Success follow user!");
              } else {
                await unfollowUser(user.id, userProfile.id);
                toastStore.success("Success unfollow user!");
              }

              await refreshData(page, activeTab);
            }}
            disabled={false}
          >
            {userProfile.is_followed ? "Unfollow" : "Follow"}
          </button>
        {/if}
      </div>
    </div>

    <div class="tabs">
      <button
        type="button"
        class="tab-button"
        class:active={activeTab === "posts"}
        on:click={() => {
          activeTab = "posts";
          page = 1;
          refreshData(page, activeTab);
        }}
      >
        Posts
      </button>
      <button
        type="button"
        class="tab-button"
        class:active={activeTab === "replies"}
        on:click={() => {
          activeTab = "replies";
          page = 1;
          refreshData(page, activeTab);
        }}
      >
        Replies
      </button>
      {#if isUser}
        <button
          type="button"
          class="tab-button"
          class:active={activeTab === "likes"}
          on:click={() => {
            activeTab = "likes";
            page = 1;
            refreshData(page, activeTab);
          }}
        >
          Likes
        </button>
      {/if}
      <button
        type="button"
        class="tab-button"
        class:active={activeTab === "media"}
        on:click={() => {
          activeTab = "media";
          page = 1;
          refreshData(page, activeTab);
        }}
      >
        Media
      </button>
    </div>

    {#if activeTab === "media"}
      <div class="media-section">
        <h2 class="media-title">Media</h2>
        <div class="media-grid">
          {#each threadMedia as mediaItem (mediaItem.id)}
            <div
              role="button"
              class="media-item"
              on:click={() => navigate(`/thread-detail/${mediaItem.thread_id}`)}
              on:keydown={(event) => {
                if (event.key === "Enter" || event.key === " ") {
                  navigate(`/thread-detail/${mediaItem.thread_id}`);
                }
              }}
              tabindex="0"
            >
              {#if mediaItem.file.endsWith(".gif") || mediaItem.file.endsWith(".jpg") || mediaItem.file.endsWith(".png") || mediaItem.file.endsWith(".jpeg")}
                <img
                  src={mediaItem.file}
                  alt="Media content"
                  class="media-content"
                />
              {:else}
                <video src={mediaItem.file} controls class="media-content">
                  <track
                    kind="captions"
                    src="path_to_captions.vtt"
                    srclang="en"
                    label="English"
                  />
                </video>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {:else}
      <ThreadFeed {threads} {user} />
    {/if}
  </main>
  <!-- Profile Picture Modal -->
  {#if showProfileModal}
    <div
      class="modal-overlay"
      on:click={handleModalClick}
      on:keydown={handleKeydown}
      role="dialog"
      aria-modal="true"
      aria-labelledby="profile-modal-title"
      tabindex="-1"
    >
      <div class="modal-content">
        <div class="modal-header">
          <h2 id="profile-modal-title" class="modal-title">
            {userProfile?.name}'s Profile Picture
          </h2>
          <button
            type="button"
            class="modal-close-button"
            on:click={closeProfileModal}
            aria-label="Close modal"
          >
            <svg
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M18 6L6 18"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M6 6L18 18"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
        </div>
        <div class="modal-image-container">
          <img
            src={userProfile?.profile_picture}
            alt="{userProfile?.name} profile picture"
            class="modal-image"
          />
        </div>
      </div>
    </div>
  {/if}
{/if}

<style>
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    margin: 0 auto;
    overflow-y: auto;
  }

  .profile-header {
    display: flex;
    flex-direction: column;
    width: 100%;
  }

  .banner-container {
    position: relative;
    width: 100%;
    height: 180px;
  }

  .profile-banner {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .back-button {
    position: absolute;
    top: 16px;
    left: 16px;
    background: rgba(0, 0, 0, 0.5);
    border: none;
    border-radius: 50%;
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .profile-info {
    display: flex;
    padding: 16px;
    position: relative;
  }

  .profile-logo {
    width: 72px;
    height: 72px;
    border-radius: 16px;
    border: 3px solid white;
    margin-top: -36px;
    margin-right: 16px;
    background-color: white;
    object-fit: cover;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .profile-details {
    flex: 1;
  }

  .profile-username {
    background-color: var(--hover-background);
    color: var(--text-tertiary);
    font-size: 14px;
    padding: 2px 8px;
    border-radius: 9999px;
  }

  .profile-name {
    font-size: 20px;
    font-weight: 700;
    margin: 0 0 4px 0;
  }

  .profile-bio {
    font-size: 14px;
    color: var(--text-secondary);
    margin: 0 0 12px 0;
  }

  .profile-joined {
    font-size: 14px;
    color: var(--text-secondary);
    margin: 0;
  }

  .follow-button {
    background-color: var(--button-primary);
    color: white;
    border: none;
    border-radius: 24px;
    padding: 8px 20px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
    height: 30px;
  }

  .tabs {
    display: flex;
    border-bottom: 1px solid var(--border-color);
  }

  .tab-button {
    flex: 1;
    background: none;
    border: none;
    color: var(--text-tertiary);
    font-size: 15px;
    font-weight: 500;
    padding: 16px;
    cursor: pointer;
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

  .media-section {
    padding: 16px;
    background-color: var(--card-background);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 16px;
  }

  .media-title {
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 12px;
  }

  .media-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
  }

  .media-item {
    cursor: pointer;
    transition: transform 0.2s;
  }

  .media-item:hover {
    transform: scale(1.05);
  }
  .media-content {
    width: 100%;
    height: auto;
    border-radius: 8px;
    object-fit: cover;
  }

  .profile-logo-button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    transition: transform 0.2s ease;
  }

  .profile-logo-button:hover {
    transform: scale(1.02);
  }

  .profile-logo-button:focus {
    outline: 2px solid var(--button-primary);
    outline-offset: 2px;
    border-radius: 16px;
  }

  /* Modal Styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
    animation: fadeIn 0.2s ease;
  }

  .modal-content {
    background-color: var(--card-background);
    border-radius: 12px;
    max-width: 90vw;
    max-height: 90vh;
    overflow: hidden;
    box-shadow:
      0 20px 25px -5px rgba(0, 0, 0, 0.1),
      0 10px 10px -5px rgba(0, 0, 0, 0.04);
    animation: scaleIn 0.2s ease;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-color);
    background-color: var(--background-secondary);
  }

  .modal-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .modal-close-button {
    background: none;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    padding: 8px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .modal-close-button:hover {
    background-color: var(--hover-background);
    color: var(--text-primary);
  }

  .modal-close-button:focus {
    outline: 2px solid var(--button-primary);
    outline-offset: 2px;
  }

  .modal-image-container {
    padding: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--card-background);
  }

  .modal-image {
    max-width: 100%;
    max-height: 70vh;
    width: auto;
    height: auto;
    border-radius: 12px;
    object-fit: contain;
    box-shadow:
      0 4px 6px -1px rgba(0, 0, 0, 0.1),
      0 2px 4px -1px rgba(0, 0, 0, 0.06);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes scaleIn {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  /* Mobile responsive styles */
  @media (max-width: 768px) {
    .content-area {
      width: 100%;
      height: 100vh;
    }

    .profile-info {
      flex-direction: column;
      align-items: flex-start;
    }

    .profile-logo {
      width: 60px;
      height: 60px;
      margin-top: -30px;
      margin-right: 0;
      margin-bottom: 12px;
    }

    .follow-button {
      align-self: flex-end;
      margin-top: 12px;
    }

    .media-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 8px;
    }

    .modal-content {
      max-width: 95vw;
      max-height: 95vh;
      margin: 16px;
    }

    .modal-header {
      padding: 12px 16px;
    }

    .modal-title {
      font-size: 16px;
    }

    .modal-image-container {
      padding: 16px;
    }

    .modal-image {
      max-height: 60vh;
    }
  }

  @media (max-width: 480px) {
    .profile-info {
      padding: 12px;
    }

    .tabs {
      font-size: 14px;
    }

    .tab-button {
      padding: 12px 8px;
      font-size: 14px;
    }

    .media-grid {
      grid-template-columns: 1fr;
    }

    .modal-header {
      padding: 10px 12px;
    }

    .modal-image-container {
      padding: 12px;
    }
  }
</style>
