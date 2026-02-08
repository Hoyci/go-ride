import { createContext, useContext, useState, ReactNode } from "react";

export type UserType = "PASSENGER" | "DRIVER";
export const userTypes = ["PASSENGER", "DRIVER"] as const;

export interface User {
  id: string;
  name: string;
  email: string;
  type: UserType;
}

interface UserContextType {
  user: User | null;
  login: (user: User) => void;
  logout: () => void;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(() => {
    const savedUser = localStorage.getItem("user_data")
    return savedUser ? JSON.parse(savedUser) : null;
  });

  const login = (u: User) => {
    console.log(u)
    setUser(u);
    localStorage.setItem("user_data", JSON.stringify(u));
  };
  
  const logout = () => {
    setUser(null);
    localStorage.removeItem("user_data");
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
  };

  return (
    <UserContext.Provider value={{ user, login, logout }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const ctx = useContext(UserContext);
  if (!ctx) throw new Error("useUser must be used within UserProvider");
  return ctx;
};
