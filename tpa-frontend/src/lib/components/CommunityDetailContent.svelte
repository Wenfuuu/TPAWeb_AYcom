<script lang="ts">
  import { onMount } from "svelte";
  import {
    acceptJoinRequest,
    fetchCommunityMedia,
    rejectJoinRequest,
    requestJoinCommunity,
    type Community,
    type CommunityMedia,
    type CommunityMember,
    type JoinCommunityRequest,
  } from "../model/community";
  import type { User } from "../model/user";
  import type { UserThread } from "../model/user";
  import ThreadFeed from "./ThreadFeed.svelte";
  import { navigate } from "svelte-routing";
  import {
    fetchCommunityThreadWithUser,
    fetchTopCommunityThreadWithUser,
    type Thread,
  } from "../model/thread";

  export let community: Community;
  export let user: User | null;
  export let threads: (Thread & { userData: UserThread })[] = [];
  export let refreshData: (
    page: number,
    currTab: "top" | "latest" | "media" | "about" | "manage"
  ) => Promise<void>;
  export let topMembers: CommunityMember[] = [];
  export let moderators: CommunityMember[] = [];
  export let joinRequests: JoinCommunityRequest[] = [];
  export let communityMedia: CommunityMedia[] = [];
  export let isModerator = false;

  let activeTab: "top" | "latest" | "media" | "about" | "manage" = "latest";
  let page = 1;
  let isLoading = false;
  let hasMore = true;

  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;

  function observeScroll() {
    if (observer) {
      observer.disconnect();
    }

    if (!contentAreaElement) {
      requestAnimationFrame(observeScroll);
      return;
    }

    observer = new IntersectionObserver(
      async (entries) => {
        // console.log(hasMore);
        if (entries[0].isIntersecting && !isLoading && hasMore) {
          isLoading = true;
          //   console.log("Loading more content, page:", page + 1);
          await new Promise((resolve) => setTimeout(resolve, 500));

          let newThreads;
          let newMedias;
          if (user) {
            if (activeTab === "latest") {
              newThreads = await fetchCommunityThreadWithUser(
                page,
                user.id,
                community.id
              );
            } else if (activeTab === "top") {
              newThreads = await fetchTopCommunityThreadWithUser(
                page,
                user.id,
                community.id
              );
            } else if (activeTab === "media") {
              newMedias = await fetchCommunityMedia(page, community.id);
            }
          }

          if (activeTab === "media") {
            if (!newMedias) {
              hasMore = false;
              isLoading = false;
              return;
            }

            communityMedia = [
              ...communityMedia,
              ...newMedias.filter(
                (newMedia) =>
                  !newMedias.some((t: CommunityMedia) => t.id === newMedia.id)
              ),
            ];
          } else {
            if (!newThreads) {
              hasMore = false;
              isLoading = false;
              return;
            }
            // console.log("Fetched new threads:", newThreads.length);
            threads = [
              ...threads,
              ...newThreads.filter(
                (newThread) => !threads.some((t) => t.id === newThread.id)
              ),
            ];
          }

          page += 1;

          isLoading = false;
        }
      },
      {
        root: contentAreaElement,
        rootMargin: "1px",
        threshold: 1,
      }
    );

    if (sentinel) {
      observer.observe(sentinel);
    }
  }

  async function joinCommunity() {
    if (!user) return;
    await requestJoinCommunity(user.id, community.id);
  }

  async function acceptRequest(userId: string, communityId: string) {
    // console.log("Accepting request");
    await acceptJoinRequest(userId, communityId);
    await refreshData(page, activeTab);
  }

  async function rejectRequest(userId: string, communityId: string) {
    // console.log("Rejecting request");
    await rejectJoinRequest(userId, communityId);
    await refreshData(page, activeTab);
  }

  onMount(async () => {
    observeScroll();
    await refreshData(page, activeTab);
  });
</script>

{#if user && community}
  <main class="content-area" bind:this={contentAreaElement}>
    <!-- Community Header -->
    <div class="community-header">
      <!-- Banner Image -->
      <div class="banner-container">
        <img
          src={community.banner}
          alt="{community.name} banner"
          class="community-banner"
        />
        <button
          type="button"
          class="back-button"
          aria-label="btn"
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

      <!-- Community Info -->
      <div class="community-info">
        <img
          src={community.logo}
          alt="{community.name} logo"
          class="community-logo"
        />
        <div class="community-details">
          <h1 class="community-name">{community.name}</h1>
          <span class="community-category">{community.category}</span>
          <p class="community-description">{community.description}</p>
          <div class="community-stats">
            <span class="members-count">9.8K Members</span>
          </div>
        </div>
        <button
          type="button"
          class="join-button"
          class:joined={community.is_joined}
          class:pending={community.is_pending}
          on:click={() => {
            joinCommunity();
            refreshData(page, activeTab);
          }}
          disabled={community.is_pending}
        >
          {#if community && community.is_joined}
            Joined
          {:else if community && community.is_pending}
            Pending
          {:else}
            Request to Join
          {/if}
        </button>
      </div>
    </div>

    <div class="tabs">
      <button
        type="button"
        class="tab-button"
        class:active={activeTab === "top"}
        on:click={() => {
          activeTab = "top";
          page = 1;
          refreshData(page, activeTab);
        }}
      >
        Top
      </button>
      <button
        type="button"
        class="tab-button"
        class:active={activeTab === "latest"}
        on:click={() => {
          activeTab = "latest";
          page = 1;
          refreshData(page, activeTab);
        }}
      >
        Latest
      </button>
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
      <button
        type="button"
        class="tab-button"
        class:active={activeTab === "about"}
        on:click={() => {
          activeTab = "about";
          page = 1;
          refreshData(page, activeTab);
        }}
      >
        About
      </button>
      {#if isModerator}
        <button
          type="button"
          class="tab-button"
          class:active={activeTab === "manage"}
          on:click={() => {
            activeTab = "manage";
            page = 1;
            refreshData(page, activeTab);
          }}
        >
          Manage
        </button>
      {/if}
    </div>

    {#if activeTab === "latest"}
      <ThreadFeed {threads} {user} />
      <div bind:this={sentinel} class="loading-sentinel">
        {#if isLoading}
          <div class="skeleton-loading">
            <!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
            {#each Array(3) as _, i (i)}
              <div class="skeleton-post">
                <div class="skeleton-avatar pulse"></div>
                <div class="skeleton-content">
                  <div class="skeleton-header">
                    <div class="skeleton-name pulse"></div>
                    <div class="skeleton-username pulse"></div>
                  </div>
                  <div class="skeleton-text pulse"></div>
                  <div class="skeleton-text-short pulse"></div>
                  {#if i % 2 === 0}
                    <div class="skeleton-image pulse"></div>
                  {/if}
                  <div class="skeleton-actions">
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        {:else if !hasMore && threads.length > 0}
          <p class="no-threads">No more thread available!</p>
        {/if}
      </div>
    {:else if activeTab === "top"}
      <!-- top members here -->
      <div class="top-members">
        <p class="top-members-title">Top Members</p>
        {#each topMembers as member (member.id)}
          <div class="follow-item">
            <img
              src={member.profile_picture || "/placeholder.svg"}
              alt={member.name}
              class="follow-avatar"
            />
            <div class="follow-info">
              <div class="follow-name">
                {member.name}
                {#if member.is_premium}
                  <span class="verified-badge" aria-label="Verified account">
                    <svg
                      viewBox="0 0 24 24"
                      width="16"
                      height="16"
                      fill="rgb(29, 155, 240)"
                      aria-hidden="true"
                    >
                      <path
                        d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"
                      ></path>
                    </svg>
                  </span>
                {/if}
              </div>
              <div class="follow-username">@{member.username}</div>
            </div>
            <button type="button" class="follow-button">Follow</button>
          </div>
        {/each}
      </div>

      <ThreadFeed {threads} {user} />
      <div bind:this={sentinel} class="loading-sentinel">
        {#if isLoading}
          <div class="skeleton-loading">
            <!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
            {#each Array(3) as _, i (i)}
              <div class="skeleton-post">
                <div class="skeleton-avatar pulse"></div>
                <div class="skeleton-content">
                  <div class="skeleton-header">
                    <div class="skeleton-name pulse"></div>
                    <div class="skeleton-username pulse"></div>
                  </div>
                  <div class="skeleton-text pulse"></div>
                  <div class="skeleton-text-short pulse"></div>
                  {#if i % 2 === 0}
                    <div class="skeleton-image pulse"></div>
                  {/if}
                  <div class="skeleton-actions">
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        {:else if !hasMore && threads.length > 0}
          <p class="no-threads">No more thread available!</p>
        {/if}
      </div>
    {:else if activeTab === "media"}
      <div class="media-section">
        <h2 class="media-title">Media</h2>
        {#if communityMedia.length > 0}
          <div class="media-grid">
            {#each communityMedia as mediaItem (mediaItem.id)}
              <div
                role="button"
                class="media-item"
                on:click={() =>
                  navigate(`/thread-detail/${mediaItem.thread_id}`)}
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
                    loading="lazy"
                  />
                {:else}
                  <video
                    src={mediaItem.file}
                    class="media-content"
                    preload="metadata"
                  >
                    <track
                      kind="captions"
                      src="path_to_captions.vtt"
                      srclang="en"
                      label="English"
                    />
                  </video>
                  <div class="video-overlay">
                    <svg
                      width="24"
                      height="24"
                      viewBox="0 0 24 24"
                      fill="white"
                    >
                      <path d="M8 5v14l11-7z" />
                    </svg>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {:else}
          <div class="empty-state">
            <div class="empty-icon">📸</div>
            <h3>No media yet</h3>
            <p>
              This community hasn't shared any photos or videos yet. Be the
              first to post some media!
            </p>
          </div>
        {/if}
        <div bind:this={sentinel} class="loading-sentinel">
          {#if isLoading}
            <div class="loading-spinner">
              <div class="spinner"></div>
              <p>Loading more media...</p>
            </div>
          {:else if !hasMore && communityMedia.length > 0}
            <div class="end-message">
              <span>✨ You've reached the end! ✨</span>
            </div>
          {/if}
        </div>
      </div>
    {:else if activeTab === "about"}
      <div class="about-section">
        <h2>About {community.name}</h2>

        <div class="about-grid">
          <div class="about-card">
            <div class="about-card-header">
              <span class="about-icon">👤</span>
              <h3>Creator</h3>
            </div>
            <p class="about-card-content">{community.creator_name}</p>
          </div>

          <div class="about-card">
            <div class="about-card-header">
              <span class="about-icon">📅</span>
              <h3>Created</h3>
            </div>
            <p class="about-card-content">
              {new Date(community.created_at).toLocaleDateString("en-US", {
                year: "numeric",
                month: "long",
                day: "numeric",
              })}
            </p>
          </div>
        </div>

        <div class="about-card rules-card">
          <h3>Community Rules</h3>
          <div class="rules-content">
            <p>{community.rule}</p>
          </div>
        </div>

        <div class="about-card moderators-card">
          <h3>Moderators</h3>
          <div class="moderators-list">
            {#each moderators as moderator (moderator.id)}
              <div class="moderator-item">
                <img
                  src={moderator.profile_picture || "/placeholder.svg"}
                  alt={moderator.name}
                  class="moderator-avatar"
                />
                <div class="moderator-info">
                  <span class="moderator-name">{moderator.name}</span>
                  <span class="moderator-username">@{moderator.username}</span>
                </div>
                <span class="moderator-badge">MOD</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    {:else if activeTab === "manage"}
      <div class="manage-section">
        <h2 class="manage-title">Manage Requests</h2>
        {#if joinRequests.length > 0}
          <ul class="manage-list">
            {#each joinRequests as request (request.id)}
              <li class="manage-list-item">
                <span class="request-info">
                  <strong>Username:</strong>
                  {request.username} - <strong>Name:</strong>
                  {request.name}
                </span>
                <div class="manage-actions">
                  <button
                    type="button"
                    class="accept-button"
                    on:click={() =>
                      acceptRequest(request.user_id, community.id)}
                    >Accept</button
                  >
                  <button
                    type="button"
                    class="reject-button"
                    on:click={() =>
                      rejectRequest(request.user_id, community.id)}
                    >Reject</button
                  >
                </div>
              </li>
            {/each}
          </ul>
        {:else}
          <p class="no-requests">No pending requests.</p>
        {/if}
      </div>
    {/if}
  </main>
{/if}

<style>
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    margin: 0 auto; /* Center if using max-width */
    overflow-y: auto;
  }

  /* Community Header Styles */
  .community-header {
    display: flex;
    flex-direction: column;
    width: 100%;
  }

  .banner-container {
    position: relative;
    width: 100%;
    height: 180px;
  }

  .community-banner {
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

  .community-info {
    display: flex;
    padding: 16px;
    position: relative;
  }

  .community-logo {
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

  .community-details {
    flex: 1;
  }

  .community-category {
    background-color: var(--hover-background);
    color: var(--text-tertiary);
    font-size: 14px;
    padding: 2px 8px;
    border-radius: 9999px;
  }

  .community-name {
    font-size: 20px;
    font-weight: 700;
    margin: 0 0 4px 0;
  }

  .community-description {
    font-size: 14px;
    color: var(--text-secondary);
    margin: 0 0 12px 0;
  }

  .community-stats {
    display: flex;
    align-items: center;
  }

  .members-count {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .join-button {
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

  .join-button.joined {
    background-color: transparent;
    border: 1px solid var(--button-primary);
    color: var(--button-primary);
  }

  .join-button.pending {
    background-color: var(--hover-background);
    color: var(--text-secondary);
    cursor: not-allowed;
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

  /* Skeleton loading styles */
  .skeleton-loading {
    width: 90%;
    padding: 0 16px;
  }

  .skeleton-post {
    display: flex;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-color);
    animation: fadeIn 0.5s ease-in-out;
  }

  .skeleton-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: var(--hover-background);
    margin-right: 12px;
    flex-shrink: 0;
  }

  .skeleton-content {
    flex: 1;
  }

  .skeleton-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    gap: 8px;
  }

  .skeleton-name {
    height: 16px;
    width: 120px;
    background-color: var(--hover-background);
    border-radius: 4px;
  }

  .skeleton-username {
    height: 14px;
    width: 80px;
    background-color: var(--hover-background);
    border-radius: 4px;
  }

  .skeleton-text {
    height: 14px;
    width: 100%;
    background-color: var(--hover-background);
    border-radius: 4px;
    margin-bottom: 8px;
  }

  .skeleton-text-short {
    height: 14px;
    width: 70%;
    background-color: var(--hover-background);
    border-radius: 4px;
    margin-bottom: 16px;
  }

  .skeleton-image {
    height: 180px;
    width: 100%;
    background-color: var(--hover-background);
    border-radius: 16px;
    margin-bottom: 16px;
  }

  .skeleton-actions {
    display: flex;
    justify-content: space-between;
    max-width: 425px;
  }

  .skeleton-action {
    height: 18px;
    width: 36px;
    background-color: var(--hover-background);
    border-radius: 4px;
  }

  .no-threads {
    padding: 32px 16px;
    text-align: center;
    color: var(--text-tertiary);
  }

  .follow-item {
    display: flex;
    padding: 12px 0;
    position: relative;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .follow-item:hover {
    background-color: var(--hover-background);
  }

  .follow-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    margin-right: 12px;
  }

  .follow-info {
    flex: 1;
    min-width: 0;
  }

  .follow-name {
    display: flex;
    align-items: center;
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .verified-badge {
    margin-left: 4px;
    display: flex;
    align-items: center;
  }

  .follow-username {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .follow-button {
    background-color: var(--text-primary);
    color: var(--card-background);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .follow-button:hover {
    opacity: 0.9;
  }
  /* Responsive Design */
  @media (max-width: 768px) {
    .content-area {
      border-right: none;
      width: 100%;
      padding: 0;
    }

    .community-header {
      margin: 0;
    }

    .banner-container {
      height: 150px;
    }

    .community-info {
      padding: 12px 16px;
      flex-direction: column;
      align-items: flex-start;
    }

    .community-logo {
      width: 60px;
      height: 60px;
      margin-top: -30px;
      margin-right: 0;
      margin-bottom: 12px;
    }

    .join-button {
      align-self: flex-end;
      margin-top: -40px;
    }

    .tabs {
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
      scrollbar-width: none;
      -ms-overflow-style: none;
    }

    .tabs::-webkit-scrollbar {
      display: none;
    }

    .tab-button {
      flex: 0 0 auto;
      min-width: 80px;
      padding: 12px 16px;
      font-size: 14px;
    }

    /* Media section mobile */
    .media-section {
      padding: 16px 12px;
      margin: 8px;
    }

    .media-title {
      font-size: 20px;
      margin-bottom: 16px;
    }

    .media-grid {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      gap: 12px;
    } /* About section mobile */
    .about-section {
      padding: 16px;
      margin: 8px;
    }

    .about-section h2 {
      font-size: 20px;
    }

    .about-grid {
      grid-template-columns: 1fr;
      gap: 12px;
    }

    .about-card {
      padding: 16px;
    }

    .about-card h3 {
      font-size: 15px;
    }

    .about-card-content {
      font-size: 14px;
    }

    .about-icon {
      font-size: 20px;
    }

    .moderator-item {
      padding: 10px 12px;
    }

    .moderator-avatar {
      width: 36px;
      height: 36px;
    }

    .moderator-name {
      font-size: 13px;
    }

    .moderator-username {
      font-size: 12px;
    }

    /* Manage section mobile */
    .manage-section {
      padding: 16px;
      margin: 8px;
    }

    .manage-title {
      font-size: 20px;
    }

    .manage-list-item {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;
      padding: 16px;
    }

    .request-info {
      font-size: 14px;
      text-align: center;
    }

    .manage-actions {
      justify-content: center;
      gap: 8px;
    }

    .accept-button,
    .reject-button {
      flex: 1;
      min-width: 0;
      padding: 12px 16px;
    }

    .top-members {
      padding: 12px;
      margin: 8px;
    }

    .follow-item {
      padding: 8px;
    }

    .follow-avatar {
      width: 40px;
      height: 40px;
    }

    .follow-name {
      font-size: 14px;
    }

    .follow-username {
      font-size: 13px;
    }

    .follow-button {
      padding: 6px 12px;
      font-size: 12px;
    }
  }

  @media (max-width: 480px) {
    .media-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 8px;
    }

    .media-title {
      font-size: 18px;
    }

    .about-section h2,
    .manage-title {
      font-size: 18px;
    }

    .manage-list-item {
      padding: 12px;
    }

    .accept-button,
    .reject-button {
      padding: 10px 12px;
      font-size: 13px;
    }

    .community-name {
      font-size: 18px;
    }

    .community-description {
      font-size: 13px;
    }
  }

  @media (max-width: 320px) {
    .media-grid {
      grid-template-columns: 1fr;
    }

    .tabs {
      font-size: 12px;
    }

    .tab-button {
      min-width: 70px;
      padding: 10px 12px;
    }
  }

  /* Large screen optimizations */
  @media (min-width: 1200px) {
    .content-area {
      max-width: 800px;
    }

    .media-grid {
      grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
      gap: 20px;
    }

    .about-section,
    .manage-section {
      padding: 32px;
      margin: 20px;
    }

    .media-section {
      padding: 32px;
      margin: 20px;
    }
  }

  .pulse {
    animation: pulse 1.5s infinite ease-in-out;
  }

  @keyframes pulse {
    0% {
      opacity: 0.6;
    }
    50% {
      opacity: 0.3;
    }
    100% {
      opacity: 0.6;
    }
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  /* Updated styling for the top members section */
  .top-members {
    padding: 16px;
    background-color: var(--card-background);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 16px;
  }

  .top-members-title {
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 12px;
  }

  .follow-item {
    display: flex;
    align-items: center;
    padding: 12px;
    border-radius: 8px;
    transition: background-color 0.2s;
    cursor: pointer;
    margin-bottom: 8px;
  }

  .follow-item:hover {
    background-color: var(--hover-background);
  }

  .follow-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    margin-right: 12px;
    object-fit: cover;
  }

  .follow-info {
    flex: 1;
    min-width: 0;
  }

  .follow-name {
    display: flex;
    align-items: center;
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .verified-badge {
    margin-left: 4px;
    display: flex;
    align-items: center;
  }

  .follow-username {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .follow-button {
    background-color: var(--button-primary);
    color: white;
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .follow-button:hover {
    background-color: var(--button-primary-hover);
  }
  /* About Section Styles */
  .about-section {
    padding: 24px;
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 4px 12px var(--shadow-color);
    margin: 16px;
    border: 1px solid var(--border-color);
  }

  .about-section h2 {
    font-size: 24px;
    font-weight: 800;
    color: var(--text-primary);
    margin: 0 0 20px 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .about-section h2::before {
    content: "ℹ️";
    font-size: 20px;
  }

  .about-section h3 {
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 24px 0 12px 0;
    border-bottom: 2px solid var(--button-primary);
    padding-bottom: 4px;
    display: inline-block;
  }

  .about-section p {
    font-size: 15px;
    line-height: 1.6;
    color: var(--text-secondary);
    margin: 0 0 16px 0;
  }

  .about-section strong {
    color: var(--text-primary);
    font-weight: 600;
  }

  .about-section ul {
    list-style: none;
    padding: 0;
    margin: 12px 0;
  }

  .about-section li {
    background-color: var(--hover-background);
    padding: 12px 16px;
    margin: 8px 0;
    border-radius: 12px;
    border-left: 4px solid var(--button-primary);
    font-size: 14px;
    color: var(--text-secondary);
    transition: all 0.2s ease;
  }

  .about-section li:hover {
    background-color: var(--search-background);
    transform: translateX(4px);
  }

  /* Manage Section Styles */
  .manage-section {
    padding: 24px;
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 4px 12px var(--shadow-color);
    margin: 16px;
    border: 1px solid var(--border-color);
  }

  .manage-title {
    font-size: 24px;
    font-weight: 800;
    color: var(--text-primary);
    margin: 0 0 20px 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .manage-title::before {
    content: "⚙️";
    font-size: 20px;
  }

  .manage-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .manage-list-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    margin: 12px 0;
    background-color: var(--hover-background);
    border-radius: 12px;
    border: 1px solid var(--border-color);
    transition: all 0.3s ease;
  }

  .manage-list-item:hover {
    background-color: var(--search-background);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px var(--shadow-color);
  }

  .request-info {
    font-size: 15px;
    color: var(--text-primary);
    font-weight: 500;
  }

  .request-info strong {
    color: var(--button-primary);
    font-weight: 600;
  }

  .manage-actions {
    display: flex;
    gap: 12px;
  }

  .accept-button,
  .reject-button {
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 600;
    border: none;
    border-radius: 25px;
    cursor: pointer;
    transition: all 0.3s ease;
    min-width: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
  }

  .accept-button {
    background: linear-gradient(135deg, #28a745, #20c997);
    color: white;
    box-shadow: 0 2px 8px rgba(40, 167, 69, 0.3);
  }

  .accept-button:hover {
    background: linear-gradient(135deg, #218838, #1abc9c);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(40, 167, 69, 0.4);
  }

  .accept-button::before {
    content: "✓";
    font-weight: bold;
  }

  .reject-button {
    background: linear-gradient(135deg, #dc3545, #e74c3c);
    color: white;
    box-shadow: 0 2px 8px rgba(220, 53, 69, 0.3);
  }

  .reject-button:hover {
    background: linear-gradient(135deg, #c82333, #c0392b);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(220, 53, 69, 0.4);
  }

  .reject-button::before {
    content: "✗";
    font-weight: bold;
  }

  .no-requests {
    font-size: 16px;
    color: var(--text-tertiary);
    text-align: center;
    padding: 40px;
    background-color: var(--hover-background);
    border-radius: 12px;
    border: 2px dashed var(--border-color);
    margin-top: 20px;
  }

  .no-requests::before {
    content: "📝";
    display: block;
    font-size: 48px;
    margin-bottom: 12px;
  }

  /* Media Section Styles */
  .media-section {
    padding: 24px;
    margin: 16px;
  }

  .media-title {
    font-size: 24px;
    font-weight: 800;
    color: var(--text-primary);
    margin: 0 0 24px 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .media-title::before {
    content: "🎬";
    font-size: 20px;
  }

  .media-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  .media-item {
    position: relative;
    cursor: pointer;
    border-radius: 16px;
    overflow: hidden;
    background-color: var(--hover-background);
    transition: all 0.3s ease;
    aspect-ratio: 1;
    border: 2px solid transparent;
  }

  .media-item:hover {
    transform: translateY(-4px) scale(1.02);
    box-shadow: 0 8px 25px var(--shadow-color);
    border-color: var(--button-primary);
  }

  .media-item::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(
      135deg,
      rgba(29, 155, 240, 0.1),
      rgba(29, 155, 240, 0.3)
    );
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .media-item:hover::after {
    opacity: 1;
  }

  .media-content {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .media-item:hover .media-content {
    transform: scale(1.1);
  }

  /* Video overlay for video thumbnails */
  .video-overlay {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: rgba(0, 0, 0, 0.6);
    border-radius: 50%;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    backdrop-filter: blur(4px);
  }

  .media-item:hover .video-overlay {
    background: rgba(29, 155, 240, 0.8);
    transform: translate(-50%, -50%) scale(1.1);
  }

  /* Empty state styling */
  .empty-state {
    text-align: center;
    padding: 60px 20px;
    background: linear-gradient(
      135deg,
      var(--hover-background),
      var(--search-background)
    );
    border-radius: 20px;
    border: 2px dashed var(--border-color);
    margin: 20px 0;
  }

  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
    opacity: 0.8;
  }

  .empty-state h3 {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 12px 0;
  }

  .empty-state p {
    font-size: 15px;
    color: var(--text-secondary);
    line-height: 1.5;
    max-width: 400px;
    margin: 0 auto;
  }

  /* Loading spinner */
  .loading-spinner {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 30px;
    gap: 16px;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--spinner-border);
    border-top: 3px solid var(--spinner-border-top);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .loading-spinner p {
    color: var(--text-secondary);
    font-size: 14px;
    margin: 0;
  }

  /* End message styling */
  .end-message {
    text-align: center;
    padding: 24px;
    background: linear-gradient(
      135deg,
      var(--button-primary),
      var(--button-primary-hover)
    );
    border-radius: 20px;
    margin: 20px 0;
    color: white;
    font-weight: 600;
    font-size: 16px;
    box-shadow: 0 4px 12px rgba(29, 155, 240, 0.3);
  }

  /* About section enhanced styling */
  .about-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  .about-card {
    background: linear-gradient(
      135deg,
      var(--card-background),
      var(--hover-background)
    );
    border: 1px solid var(--border-color);
    border-radius: 16px;
    padding: 20px;
    transition: all 0.3s ease;
    box-shadow: 0 2px 8px var(--shadow-color);
  }

  .about-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px var(--shadow-color);
    border-color: var(--button-primary);
  }

  .about-card-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
  }

  .about-icon {
    font-size: 24px;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  }

  .about-card h3 {
    font-size: 16px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
  }

  .about-card-content {
    font-size: 15px;
    color: var(--text-secondary);
    margin: 0;
    font-weight: 500;
  }

  .rules-card,
  .moderators-card {
    grid-column: 1 / -1;
  }

  .rules-content {
    background-color: var(--search-background);
    padding: 16px;
    border-radius: 12px;
    border-left: 4px solid var(--button-primary);
    margin-top: 12px;
  }

  .rules-content p {
    margin: 0;
    line-height: 1.6;
    color: var(--text-primary);
  }

  .moderators-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 16px;
  }

  .moderator-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background-color: var(--search-background);
    border-radius: 12px;
    transition: all 0.2s ease;
    border: 1px solid transparent;
  }

  .moderator-item:hover {
    background-color: var(--hover-background);
    border-color: var(--button-primary);
    transform: translateX(4px);
  }

  .moderator-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid var(--button-primary);
  }

  .moderator-info {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .moderator-name {
    font-weight: 600;
    color: var(--text-primary);
    font-size: 14px;
  }

  .moderator-username {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .moderator-badge {
    background: linear-gradient(
      135deg,
      var(--button-primary),
      var(--button-primary-hover)
    );
    color: white;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.5px;
  }
</style>
