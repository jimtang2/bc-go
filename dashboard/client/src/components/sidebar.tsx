import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader, SidebarSeparator, SidebarGroup, SidebarMenu, SidebarMenuItem, SidebarMenuButton } from '@/components/ui/sidebar';
import { Blend, History, Settings } from 'lucide-react';
import clsx from "clsx";
const links = {
  alpha: { to: '/', label: 'Home', paths: ['/', '/alpha'] },
  settings: { to: '/settings', label: 'Settings' },
  history: { to: '/history', label: 'History' },
};
const AppSidebar: React.FC<{ variant?: 'sidebar' | 'floating' | 'inset' }> = ({ variant }) => {
  
  return (
    <Sidebar variant={variant}>
      <SidebarHeader>
        <div className="flex flex-col gap-2 p-2">
          <span className="text-lg font-semibold">Arbitrage</span>
        </div>
      </SidebarHeader>
      <SidebarContent className="overflow-x-hidden">
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton>
                <SidebarMenuButtonLink {...links.alpha}>
                  <Blend className="mr-2 h-4 w-4" />
                </SidebarMenuButtonLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton>
                <SidebarMenuButtonLink {...links.history}>
                  <History className="mr-2 h-4 w-4" />
                </SidebarMenuButtonLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
        <SidebarSeparator />
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton>
                <SidebarMenuButtonLink {...links.settings}>
                  <Settings className="mr-2 h-4 w-4" />
                </SidebarMenuButtonLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton>
              {/*<Github className="mr-2 h-4 w-4" />
              <span>Github</span>*/}
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
interface SidebarMenuButtonLinkProps {
  to: string;
  label: string;
  paths?: string[];
  children?: React.ReactNode;
}
function SidebarMenuButtonLink({ to, label, paths, children }: SidebarMenuButtonLinkProps) {
  const { pathname } = useLocation();
  const isActive = paths ? paths.includes(pathname) : pathname === to;
  const className = clsx([
    "flex flex-row items-center", 
    isActive ? "text-gray-100" : "text-gray-500",
  ])
  return (
    <Link key={to} to={to} className={className}>
      {children}
      <span>{label}</span>
    </Link>
    )
}
export default AppSidebar;