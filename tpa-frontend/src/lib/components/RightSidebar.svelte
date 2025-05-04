<script lang="ts">
  import { onMount } from "svelte";
  import { fetchWhoToFolllow, type User } from "../model/user";
  import { checkAuth, getAuthUser } from "../utils/auth";
  import { navigate } from "svelte-routing";

  // Sample trending topics
  const trendingTopics = [
    {
      category: "Politics · Trending",
      title: "Trump",
      posts: "1,209 posts"
    },
    {
      category: "NCT Dream · Trending",
      title: "#NCT꿈",
      posts: "300K posts"
    },
    {
      category: "Trending in Indonesia",
      title: "Pagi",
      posts: "12.5K posts"
    },
    {
      category: "Trending in Indonesia",
      title: "#teddysoekinspirasI",
      posts: "2,543 posts"
    }
  ];

  let user: User | null;
  let toFollow: User[] = [];

  onMount(async() => {
    await checkAuth(navigate);

    const usr = await getAuthUser();
    if(usr) {
      user = usr;
      toFollow = await fetchWhoToFolllow(user.id);
    }
  });
</script>

<div class="sidebar right-sidebar">
  <div class="search-container">
    <div class="search-box">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.3-4.3"></path>
      </svg>
      <input type="text" placeholder="Search" />
    </div>
  </div>

  <div class="sidebar-section premium-section">
    <h2>Subscribe to Premium</h2>
    <p>Subscribe to unlock new features and if eligible, receive a share of revenue.</p>
    <button type="button" class="subscribe-button">Subscribe</button>
  </div>

  <div class="sidebar-section trending-section">
    <h2>What's happening</h2>
    <div class="trending-topics">
      {#each trendingTopics as topic (topic.title)}
        <div class="trending-item">
          <div class="trending-item-category">{topic.category}</div>
          <div class="trending-item-title">{topic.title}</div>
          <div class="trending-item-posts">{topic.posts}</div>
        </div>
      {/each}
    </div>
  </div>

  <div class="sidebar-section follow-section">
    <h2>Who to follow</h2>
    <div class="follow-suggestions">
      {#each toFollow as user (user.id)}
        <div class="follow-item">
          <img src={user.profile_picture || "/placeholder.svg"} alt={user.name} class="follow-avatar" />
          <div class="follow-info">
            <div class="follow-name">
              {user.name}
              {#if user.is_premium}
                <span class="verified-badge">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="rgb(29, 155, 240)">
                    <path d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"></path>
                  </svg>
                </span>
              {/if}
            </div>
            <div class="follow-username">@{user.username}</div>
          </div>
          <button type="button" class="follow-button">Follow</button>
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .right-sidebar {
    width: 25%;
    padding: 0 16px;
    position: sticky;
    top: 0;
    height: 100vh;
    overflow-y: auto;
  }

  .search-container {
    position: sticky;
    top: 0;
    padding: 6px 0;
    background-color: var(--background);
    z-index: 10;
  }

  .search-box {
    display: flex;
    align-items: center;
    background-color: var(--hover-background);
    border-radius: 9999px;
    padding: 12px 16px;
    margin: 6px 0;
    transition: background-color 0.2s;
  }

  .search-box:focus-within {
    background-color: var(--card-background);
    border: 1px solid var(--button-primary);
  }

  .search-box svg {
    color: var(--text-tertiary);
    margin-right: 12px;
  }

  .search-box input {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 15px;
    width: 100%;
    outline: none;
  }

  .search-box input::placeholder {
    color: var(--text-tertiary);
  }

  .sidebar-section {
    background-color: var(--card-background, var(--hover-background));
    border-radius: 16px;
    margin: 16px 0;
    padding: 16px;
    border-color: var(--divider-color);
    border: solid 1px;
  }

  .sidebar-section h2 {
    font-size: 20px;
    font-weight: 800;
    margin: 0 0 16px 0;
    color: var(--text-primary);
  }

  .premium-section p {
    font-size: 15px;
    margin: 0 0 16px 0;
    color: var(--text-primary);
  }

  .subscribe-button {
    background-color: var(--text-primary);
    color: var(--card-background);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .subscribe-button:hover {
    opacity: 0.9;
  }

  .trending-topics {
    display: flex;
    flex-direction: column;
  }

  .trending-item {
    padding: 12px 0;
    position: relative;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .trending-item:hover {
    background-color: var(--hover-background);
  }

  .trending-item-category {
    font-size: 13px;
    color: var(--text-tertiary);
    margin-bottom: 4px;
  }

  .trending-item-title {
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .trending-item-posts {
    font-size: 13px;
    color: var(--text-tertiary);
  }

  .follow-suggestions {
    display: flex;
    flex-direction: column;
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

  @media (max-width: 1024px) {
    .right-sidebar {
      display: none
    }
  }
</style>