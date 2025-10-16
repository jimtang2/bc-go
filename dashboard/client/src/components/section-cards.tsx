// webhook$$bc-go;dashboard/client/src/components/section-cards.tsx;grok$$
import React from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardFooter,
  CardTitle,
} from '@/components/ui/card';
import { PortfolioChart, AverageMarginChart } from './section-charts';

const SectionCards: React.FC = () => {
  return (
    <div className="*:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:shadow-xs lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Portfolio</CardDescription>
        </CardHeader>
        <CardContent className="px-2 pt-0 sm:px-6">
          <PortfolioChart />
        </CardContent>
        {/* <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Placeholder Data
          </div>
          <div className="text-muted-foreground">Static Info</div>
        </CardFooter> */}
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Average Margin</CardDescription>
        </CardHeader>
        <CardContent className="px-2 pt-0 sm:px-6">
          <AverageMarginChart />
        </CardContent>
        {/* <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Placeholder Data
          </div>
          <div className="text-muted-foreground">Static Info</div>
        </CardFooter> */}
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Trades</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            Placeholder
          </CardTitle>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Placeholder Data
          </div>
          <div className="text-muted-foreground">Static Info</div>
        </CardFooter>
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Growth Rate</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            Placeholder
          </CardTitle>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Placeholder Data
          </div>
          <div className="text-muted-foreground">Static Info</div>
        </CardFooter>
      </Card>
    </div>
  );
};

export default SectionCards;