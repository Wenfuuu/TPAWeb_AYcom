<script lang="ts">
  import { onMount } from "svelte";
  import {
    requestJoinCommunity,
    type Community,
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
            }
          }

          if (!newThreads) {
            hasMore = false;
            isLoading = false;
            return;
          }
          console.log("Fetched new threads:", newThreads.length);
          threads = [
            ...threads,
            ...newThreads.filter(
              (newThread) => !threads.some((t) => t.id === newThread.id)
            ),
          ];
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

  async function acceptRequest(requestId: string) {
    console.log(`Accepted request for ID: ${requestId}`);

  }

  async function rejectRequest(requestId: string) {
    console.log(`Rejected request for ID: ${requestId}`);

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
          on:click={joinCommunity}
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
    </div>

    {#if activeTab === "latest"}
      <ThreadFeed {threads} {user} />
      <div bind:this={sentinel} class="loading-sentinel">
        {#if isLoading}
          <div class="skeleton-loading">
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
      <p>media here</p>
    {:else if activeTab === "about"}
      <div class="about-section">
        <h2>About {community.name}</h2>
        <p><strong>Creator:</strong> {community.creator_name}</p>
        <p>
          <strong>Created on:</strong>
          {new Date(community.created_at).toLocaleDateString()}
        </p>
        <h3>Community Rules</h3>
        <p>
          {community.rule}
        </p>
        <h3>Moderators</h3>
        <ul>
          {#each moderators as moderator (moderator.id)}
            <li>{moderator.name} (@{moderator.username})</li>
          {/each}
        </ul>
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
                  <button type="button"
                    class="accept-button"
                    on:click={() => acceptRequest(request.id)}
                    >Accept</button
                  >
                  <button type="button"
                    class="reject-button"
                    on:click={() => rejectRequest(request.id)}
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
    margin-bottom: 2px;
  }

  .verified-badge {
    margin-left: 4px;
    display: flex;
    align-items: center;
  }

  .follow-username {
    font-size: 14px;
    color: var(--text-tertiary);
    margin-bottom: 2px;
  }

  .follow-button {
    background-color: var(--text-primary);
    color: var(--card-background);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    height: fit-content;
    align-self: center;
    transition: background-color 0.2s;
  }

  .follow-button:hover {
    opacity: 0.9;
  }

  @media (max-width: 768px) {
    .content-area {
      border-right: none; /* Remove border on mobile */
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

  /* Updated styling for the about section */
  .about-section {
    padding: 16px;
    background-color: var(--card-background);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 16px;
  }

  .about-title {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 12px;
  }

  .about-detail {
    font-size: 14px;
    color: var(--text-secondary);
    margin-bottom: 8px;
  }

  .about-subtitle {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin-top: 16px;
    margin-bottom: 8px;
  }

  .about-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .about-list-item {
    font-size: 14px;
    color: var(--text-secondary);
    margin-bottom: 4px;
  }

  /* Updated styling for the manage section */
  .manage-section {
    padding: 16px;
    background-color: var(--card-background);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 16px;
  }

  .manage-title {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 12px;
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
    padding: 12px;
    border-bottom: 1px solid var(--border-color);
  }

  .request-info {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .manage-actions {
    display: flex;
    gap: 8px;
  }

  .accept-button,
  .reject-button {
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

  .reject-button {
    background-color: var(--text-primary);
    color: var(--card-background);
  }

  .accept-button:hover {
    background-color: var(--button-primary-hover);
  }

  .reject-button:hover {
    opacity: 0.9;
  }

  .no-requests {
    font-size: 14px;
    color: var(--text-secondary);
    text-align: center;
    margin-top: 16px;
  }
</style>
