<script lang="ts">
  import { onMount } from "svelte";
  import { requestJoinCommunity, type Community } from "../model/community";
  import { fetchCommunityThreadWithUser, type Thread } from "../model/thread";
  import type { User } from "../model/user";
  import type { UserThread } from "../model/user";
  import ThreadFeed from "./ThreadFeed.svelte";
  import { navigate } from "svelte-routing";

  export let community: Community;
  export let user: User | null;
  export let threads: (Thread & { userData: UserThread })[] = [];
  export let refreshData: (
    page: number,
    currTab: "top" | "latest" | "media" | "about" | "manage"
  ) => Promise<void>;

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
            newThreads = await fetchCommunityThreadWithUser(
              page,
              user.id,
              community.id
            );
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

  async function joinCommunity(){
    if(!user) return;
    await requestJoinCommunity(user.id, community.id);
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
      <p>top here</p>
    {:else if activeTab === "media"}
      <p>media here</p>
    {:else if activeTab === "about"}
      <p>about here</p>
    {:else if activeTab === "manage"}
      <p>manage here</p>
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
</style>
