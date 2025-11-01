import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader, SidebarSeparator, SidebarGroup, SidebarMenu, SidebarMenuItem, SidebarMenuButton } from '@/components/ui/sidebar';
import { Scale, Blend, ArrowLeftRight, Settings } from 'lucide-react';
import clsx from "clsx";
function SidebarMenuButtonLink({to, children}: {to: string; children?: React.ReactNode;}) {
  const { pathname } = useLocation();
  const className = clsx(["flex flex-row items-center w-full gap-2", pathname === to ? "text-gray-100" : "text-gray-500"])
  return <Link to={to} className={className}>{children}</Link>
}
export default () => {
  return (
    <Sidebar variant="inset">
      <SidebarHeader>
        <div className="flex flex-row items-center gap-3 p-2">
          <Scale className="h-5 w-5" />
          <span className="text-lg font-semibold">Arbitrage</span>
        </div>
      </SidebarHeader>
      <SidebarContent className="overflow-x-hidden">
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton>
                <SidebarMenuButtonLink to="/">
                  <ArrowLeftRight className="h-4 w-4" />
                  <span>Spreads</span>
                </SidebarMenuButtonLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton>
                <SidebarMenuButtonLink to="/alpha">
                  <Blend className="h-4 w-4" />
                  <span>Alpha</span>
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
                <SidebarMenuButtonLink to="/settings">
                  <Settings className="h-4 w-4" />
                  <span>Settings</span>
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
};