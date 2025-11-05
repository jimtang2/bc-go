import React from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { SidebarInset, SidebarProvider, SidebarSeparator, SidebarTrigger, Sidebar, SidebarContent, SidebarFooter, SidebarHeader, SidebarGroup, SidebarMenu, SidebarMenuItem, SidebarMenuButton } from '@/components/ui/sidebar';
import { Scale, Blend, ArrowLeftRight } from 'lucide-react';
import clsx from "clsx";

const Layout: React.FC = () => {
  const sidebarProviderProps = {
    style: {
      '--sidebar-width': 'calc(var(--spacing) * 50)',
      '--header-height': 'calc(var(--spacing) * 12)',
    } as React.CSSProperties,
    defaultOpen: true,
  }
  return (
    <SidebarProvider {...sidebarProviderProps}>
      <AppSidebar />
      <SidebarInset>
      	<header className="flex h-[var(--header-height)] shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-[var(--header-height)]">
      	  <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
      	    <SidebarTrigger className="-ml-1 cursor-pointer" />
      	    <SidebarSeparator orientation="vertical" className="mx-2 data-[orientation=vertical]:h-4" />
      	    <h1 className="text-base font-medium flex-grow-1">Arbitrage</h1>
      	  </div>
      	</header>
        <div className="flex flex-1 flex-col">
          <div className="@container/main flex flex-1 flex-col gap-2">
            <div className="flex flex-col gap-4 py-4 md:gap-4 md:py-4">
              <Outlet />
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
};

export default Layout

interface sidebarLinkProps {
  to:     string;
  label:  string;
  icon:   React.ReactElement;
}

function AppSidebar() {
  const { pathname } = useLocation();
  const sidebarLinks: sidebarLinkProps[] = [
    {
      to: "/",
      label: "Spreads",
      icon: <ArrowLeftRight className="h-4 w-4" />,
    },
    {
      to: "/alpha/",
      label: "Alpha",
      icon: <Blend className="h-4 w-4" />,
    },
  ]
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
            {sidebarLinks.map(({to, label, icon}) => {
              const className = clsx(["flex flex-row items-center w-full gap-2", pathname === to || pathname === to + "/" ? "text-gray-100" : "text-gray-500"])
              return (
                <SidebarMenuItem key={to}>
                  <SidebarMenuButton>
                    <Link to={to} className={className}>
                      {icon}<span>{label}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                )
            })}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter></SidebarFooter>
    </Sidebar>
  );
};