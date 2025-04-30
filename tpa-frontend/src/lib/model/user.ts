export type User = {
    id: string
    name: string
    username: string
    email: string
    password: string
    gender: string
    profilePicture: string
    banner: string
    dateOfBirth: string
    isActive: boolean
    isBanned: boolean
    role: string
    isSubscribed: boolean
    isPremium: boolean
    bio: string
    createdAt: string
}

export type UserThread = {
    name: string;
    username: string;
    verified: boolean;
    profile: string;
};