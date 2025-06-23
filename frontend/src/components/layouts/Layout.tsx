import React from "react";

interface AuthLayoutProps {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
  backgroundUrl?: string;
}

const Layout: React.FC<AuthLayoutProps> = ({
  title,
  subtitle,
  children,
  backgroundUrl = "/images/bg-mountain.avif", // mountain bg
}) => {
  return (
    <div className="min-h-screen flex items-center justify-center relative overflow-hidden">
      {/* Scenic Background */}
      <div
        className="absolute inset-0 bg-cover bg-center z-0"
        style={{
          backgroundImage: `url('${backgroundUrl}')`,
        }}
      >
        <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />
      </div>

      {/* Glassmorphic Auth Card */}
      <div className="relative z-10 w-full max-w-md bg-white/10 backdrop-blur-lg rounded-xl shadow-2xl p-8 text-white">
        <div className="mb-6 text-center">
          <h1 className="text-3xl font-bold">{title}</h1>
          {subtitle && <p className="text-sm mt-2 text-gray-200">{subtitle}</p>}
        </div>
        <div>{children}</div>
      </div>
    </div>
  );
};

export default Layout;
