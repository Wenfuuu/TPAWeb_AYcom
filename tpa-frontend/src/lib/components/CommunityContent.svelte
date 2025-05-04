<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { User } from "../model/user";
  import { fetchCommunityCategory, type Community, type CommunityCategory } from "../model/community";

  // Props
  export let communities: Community[] = [];
  export let user: User | null = null;
  export let refreshData: (page: number, search?: string) => Promise<void>;
  export let changeLimit: (limit: 25 | 30 | 35) => void;

  // State variables
  let searchQuery = "";
  let selectedCategory = "all";
  let currentPage = 1;
  let totalPages = 10; // This would typically come from an API response
  let currentLimit: 25 | 30 | 35 = 25;
  let isLoading = false;
  let activeTab: "joined" | "pending" | "available" = "available";

  // Mock data for joined and pending communities (in a real app, these would come from API)
  let joinedCommunities: Community[] = [];
  let pendingCommunities: Community[] = [];

  // Categories for filtering
  let categories: CommunityCategory[] = [];

  async function handleSearch() {
    currentPage = 1;
    isLoading = true;
    await refreshData(currentPage, searchQuery);
    isLoading = false;
  }

  async function handleCategoryChange() {
    currentPage = 1;
    isLoading = true;
    // In a real app, you would include the category in the API call
    await refreshData(currentPage, searchQuery);
    isLoading = false;
  }

  async function goToPage(page: number) {
    if (page < 1 || page > totalPages) return;
    currentPage = page;
    isLoading = true;
    await refreshData(currentPage, searchQuery);
    isLoading = false;
  }

  async function handleLimitChange(limit: 25 | 30 | 35) {
    currentLimit = limit;
    changeLimit(limit);
    currentPage = 1;
    isLoading = true;
    await refreshData(currentPage, searchQuery);
    isLoading = false;
  }

  // Function to request to join a community
  async function requestToJoin(communityId: string) {
    if (!user) return;

    try {
      // In a real app, you would make an API call here
      console.log(`Requesting to join community ${communityId}`);

      // Mock implementation: move from available to pending
      const community = communities.find(c => c.id === communityId);
      if (community) {
        pendingCommunities = [...pendingCommunities, community];
        // Refresh the data to update the UI
        await refreshData(currentPage, searchQuery);
      }
    } catch (error) {
      console.error("Error requesting to join community:", error);
    }
  }

  // Function to navigate to community detail page
  function navigateToCommunity(communityId: string) {
    navigate(`/community/${communityId}`);
  }

  // Function to navigate to create community page
  function navigateToCreateCommunity() {
    navigate("/create-community");
  }

  // Function to get displayed communities based on active tab
  function getDisplayedCommunities() {
    switch (activeTab) {
      case "joined":
        return joinedCommunities;
      case "pending":
        return pendingCommunities;
      case "available":
        return communities.filter(community => {
          // Filter by category if not "all"
          if (selectedCategory !== "all" && community.category !== selectedCategory) {
            return false;
          }
          return true;
        });
      default:
        return [];
    }
  }

  // Generate page numbers for pagination
  function getPageNumbers(): number[] {
    const pages = [];
    const maxVisiblePages = 5;

    if (totalPages <= maxVisiblePages) {
      // If total pages is less than max visible, show all pages
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Always include first page
      pages.push(1);

      // Calculate start and end of visible pages
      let start = Math.max(2, currentPage - Math.floor(maxVisiblePages / 2));
      const end = Math.min(totalPages - 1, start + maxVisiblePages - 3);

      // Adjust start if end is maxed out
      if (end === totalPages - 1) {
        start = Math.max(2, end - (maxVisiblePages - 3));
      }

      // Add visible pages
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      // Always include last page
      pages.push(totalPages);
    }

    return pages;
  }

  onMount(async () => {
    isLoading = true;

    categories = await fetchCommunityCategory();

    // In a real app, you would fetch totalPages from API response
    // Example: const response = await fetchCommunities(1, searchQuery);
    // totalPages = response.totalPages;

    // fetch joined and pending communities later

    isLoading = false;
  });
</script>

<main class="content-area">
  <div class="header">
    <h1>Communities</h1>
    <button type="button" class="create-button" on:click={navigateToCreateCommunity}>
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
      Create Community
    </button>
  </div>

  <div class="search-filter-container">
    <div class="search-container">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.3-4.3"></path>
      </svg>
      <input
        type="text"
        placeholder="Search communities"
        bind:value={searchQuery}
        on:keydown={(e) => e.key === "Enter" && handleSearch()}
      />
      {#if searchQuery}
        <button type="button" class="clear-search" aria-label="clear button" on:click={() => { searchQuery = ""; handleSearch(); }}>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      {/if}
    </div>

    <div class="filter-container">
      <div class="custom-select">
        <select bind:value={selectedCategory} on:change={handleCategoryChange}>
          <option value="all">All Categories</option>
          {#each categories as category (category.id)}
            <option value={category.id}>{category.name}</option>
          {/each}
        </select>
        <div class="select-arrow">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m6 9 6 6 6-6"></path>
          </svg>
        </div>
      </div>
    </div>
  </div>

  <div class="tabs">
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "joined"}
      on:click={() => activeTab = "joined"}
    >
      Joined
      <span class="count">{joinedCommunities.length}</span>
    </button>
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "pending"}
      on:click={() => activeTab = "pending"}
    >
      Pending
      <span class="count">{pendingCommunities.length}</span>
    </button>
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "available"}
      on:click={() => activeTab = "available"}
    >
      Available
    </button>
  </div>

  {#if isLoading}
    <div class="loading-container">
      <div class="spinner"></div>
      <p>Loading communities...</p>
    </div>
  {:else}
    <div class="communities-container">
      {#if getDisplayedCommunities().length === 0}
        <div class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
          <p>No communities found</p>
          {#if activeTab === "available"}
            <button type="button" class="primary-button" on:click={navigateToCreateCommunity}>Create a Community</button>
          {/if}
        </div>
      {:else}
        <div class="communities-grid">
          {#each getDisplayedCommunities() as community (community.id)}
            <div class="community-card" role="button" tabindex="0" on:click={() => navigateToCommunity(community.id)}
                on:keydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                    navigateToCommunity(community.id);
                    e.preventDefault(); // Prevent default scrolling behavior for the Space key
                    }
                }}>
              <!-- <div class="community-banner"></div> -->
               <img
                src={community.banner || "/placeholder-banner.svg"}
                alt={`${community.name} banner`}
                class="community-banner"
              />
              <div class="community-logo-container">
                <img src={community.logo || "/placeholder.svg?height=60&width=60"} alt={community.name} class="community-logo" />
              </div>
              <div class="community-content">
                <h3 class="community-name">{community.name}</h3>
                <div class="community-category">{community.category}</div>
                <p class="community-description">{community.description}</p>

                {#if activeTab === "available"}
                  <button
                    type="button"
                    class="join-button"
                    on:click={(e) => {
                      e.stopPropagation();
                      requestToJoin(community.id);
                    }}
                  >
                    Request to Join
                  </button>
                {:else if activeTab === "pending"}
                  <div class="pending-badge">Pending Approval</div>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <div class="pagination-container">
        <div class="pagination-info">
          <span>Page {currentPage} of {totalPages}</span>
          <div class="limit-selector">
            <span>Show:</span>
            <button
              type="button"
              class="limit-button"
              class:active={currentLimit === 25}
              on:click={() => handleLimitChange(25)}
            >25</button>
            <button
              type="button"
              class="limit-button"
              class:active={currentLimit === 30}
              on:click={() => handleLimitChange(30)}
            >30</button>
            <button
              type="button"
              class="limit-button"
              class:active={currentLimit === 35}
              on:click={() => handleLimitChange(35)}
            >35</button>
          </div>
        </div>

        <div class="pagination-controls">
          <button type="button"
            aria-label="pagination button"
            class="pagination-button"
            disabled={currentPage === 1}
            on:click={() => goToPage(currentPage - 1)}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m15 18-6-6 6-6"></path>
            </svg>
          </button>

          {#each getPageNumbers() as page (page)}
            <button type="button"
              class="pagination-button"
              class:active={currentPage === page}
              on:click={() => goToPage(page)}
            >{page}</button>
          {/each}

          <button type="button"
            aria-label="pagination button"
            class="pagination-button"
            disabled={currentPage === totalPages}
            on:click={() => goToPage(currentPage + 1)}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </button>
        </div>
    </div>
  {/if}
</main>

<style>
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    margin: 0 auto;
    overflow-y: auto;
  }

  .header {
    position: sticky;
    top: 0;
    background-color: var(--card-background);
    backdrop-filter: blur(12px);
    z-index: 10;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
  }

  .header h1 {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
  }

  .create-button {
    display: flex;
    align-items: center;
    gap: 8px;
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .create-button:hover {
    background-color: var(--button-primary-hover);
  }

  .search-filter-container {
    display: flex;
    gap: 12px;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .search-container {
    position: relative;
    flex: 1;
    display: flex;
    align-items: center;
    background-color: var(--hover-background);
    border-radius: 9999px;
    padding: 0 12px;
  }

  .search-container svg {
    color: var(--text-tertiary);
    margin-right: 8px;
  }

  .search-container input {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 15px;
    padding: 10px 0;
    width: 100%;
    outline: none;
  }

  .search-container input::placeholder {
    color: var(--text-tertiary);
  }

  .clear-search {
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    padding: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .filter-container {
    min-width: 150px;
  }

  /* Improved select styling */
  .custom-select {
    position: relative;
    width: 100%;
  }

  .custom-select select {
    width: 100%;
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 9999px;
    padding: 10px 16px;
    color: var(--text-primary);
    font-size: 14px;
    outline: none;
    cursor: pointer;
    appearance: none;
    padding-right: 40px;
  }

  .custom-select select:focus {
    border-color: var(--button-primary);
    box-shadow: 0 0 0 2px rgba(var(--button-primary-rgb, 29, 155, 240), 0.2);
  }

  .custom-select select option {
    background-color: var(--card-background);
    color: var(--text-primary);
  }

  .select-arrow {
    position: absolute;
    top: 50%;
    right: 16px;
    transform: translateY(-50%);
    pointer-events: none;
    color: var(--text-tertiary);
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
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background-color: var(--button-primary);
    border-radius: 4px 4px 0 0;
  }

  .count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background-color: var(--hover-background);
    color: var(--text-tertiary);
    border-radius: 9999px;
    min-width: 20px;
    height: 20px;
    font-size: 12px;
    margin-left: 4px;
    padding: 0 6px;
  }

  .tab-button.active .count {
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
    color: var(--button-primary);
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 0;
    color: var(--text-tertiary);
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(var(--button-primary-rgb, 29, 155, 240), 0.3);
    border-radius: 50%;
    border-top-color: var(--button-primary);
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 16px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 0;
    color: var(--text-tertiary);
    text-align: center;
  }

  .empty-state svg {
    margin-bottom: 16px;
    color: var(--text-tertiary);
  }

  .empty-state p {
    font-size: 18px;
    font-weight: 700;
    margin-bottom: 16px;
    color: var(--text-primary);
  }

  .primary-button {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 10px 20px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .primary-button:hover {
    background-color: var(--button-primary-hover);
  }

  .communities-container {
    padding: 16px;
  }

  .communities-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;
  }

  .community-card {
    background-color: var(--card-background);
    border: 1px solid var(--border-color);
    border-radius: 16px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
    position: relative;
  }

  .community-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .community-banner {
    width: 100%;
    height: 100px;
    object-fit: cover;
    background-color: var(--hover-background);
  }

  .community-logo-container {
    position: absolute;
    top: 70px;
    left: 16px;
    width: 60px;
    height: 60px;
    border-radius: 12px;
    overflow: hidden;
    border: 4px solid var(--card-background);
    background-color: var(--card-background);
  }

  .community-logo {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .community-content {
    padding: 16px;
    padding-top: 36px;
  }

  .community-name {
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 4px 0;
  }

  .community-category {
    display: inline-block;
    background-color: var(--hover-background);
    color: var(--text-tertiary);
    font-size: 12px;
    font-weight: 500;
    padding: 2px 8px;
    border-radius: 9999px;
    margin-bottom: 12px;
  }

  .community-description {
    font-size: 14px;
    color: var(--text-secondary);
    margin: 0 0 16px 0;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.4;
    max-height: 4.2em;
  }

  .join-button {
    width: 100%;
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 8px 0;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .join-button:hover {
    background-color: var(--button-primary-hover);
  }

  .pending-badge {
    display: inline-block;
    width: 100%;
    text-align: center;
    background-color: var(--hover-background);
    color: var(--text-tertiary);
    font-size: 14px;
    font-weight: 500;
    padding: 8px 0;
    border-radius: 9999px;
  }

  .pagination-container {
    padding: 16px;
    border-top: 1px solid var(--border-color);
  }

  .pagination-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    color: var(--text-tertiary);
    font-size: 14px;
  }

  .limit-selector {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .limit-button {
    background: none;
    border: 1px solid var(--border-color);
    color: var(--text-tertiary);
    border-radius: 4px;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .limit-button:hover {
    background-color: var(--hover-background);
  }

  .limit-button.active {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border-color: var(--button-primary);
  }

  .pagination-controls {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
  }

  .pagination-button {
    background: none;
    border: 1px solid var(--border-color);
    color: var(--text-tertiary);
    border-radius: 4px;
    min-width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    cursor: pointer;
    transition: background-color 0.2s;
    padding: 0 8px;
  }

  .pagination-button:hover:not(:disabled) {
    background-color: var(--hover-background);
  }

  .pagination-button.active {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border-color: var(--button-primary);
  }

  .pagination-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  @media (max-width: 768px) {
    .content-area {
      width: 100%;
      border-right: none;
    }

    .communities-grid {
      grid-template-columns: 1fr;
    }

    .search-filter-container {
      flex-direction: column;
    }

    .filter-container {
      width: 100%;
    }
  }
</style>