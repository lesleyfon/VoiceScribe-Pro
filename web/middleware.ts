import { clerkMiddleware, createRouteMatcher } from "@clerk/nextjs/server";
import { NextResponse } from "next/server";

const isProtectedRoute = createRouteMatcher(["/dashboard(.*)", "/"]);
export default clerkMiddleware(async (auth, req) => {
	if (isProtectedRoute(req)) {
		const { userId } = await auth();
		if (!userId) {
			// Redirect to sign-in page if not authenticated
			return NextResponse.redirect(new URL("/auth", req.url));
		}
	}
	// Otherwise, continue as normal
	return NextResponse.next();
});
export const config = {
	match: [
		// Skip Next.js internals and all static files, unless found in search params
		"/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)",
		// Always run for API routes
		"/(api|trpc)(.*)",
	],
};
