import React, { ReactNode, createContext } from "react";

export declare interface User {
  id: string;
  email: string;
}

const initialUser: User = {
  id: "",
  email: "",
};

export const UserContext = createContext<{
  user: User;
  setUser: React.Dispatch<React.SetStateAction<User>>;
}>({
  user: initialUser,
  setUser: () => {},
});

export declare interface Session {
  user: User;
}

export declare interface UserProps {
  children?: ReactNode;
  session: Session;
}

export const UserProvider = ({ children, session }: UserProps) => {
  const sessionUser: User = { id: "", email: "" };
  sessionUser.email = session.user.email;
  const [user, setUser] = React.useState(sessionUser);

  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};
