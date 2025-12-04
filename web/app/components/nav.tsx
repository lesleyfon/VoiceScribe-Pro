"use client";

import { SignedIn, UserButton, useAuth } from "@clerk/nextjs";
import { useRouter } from "next/navigation";

export default function Nav() {
	const { isSignedIn } = useAuth();
	const router = useRouter();

	if (!isSignedIn) {
		router.push("/auth");
		return null;
	}
	return (
		<div className="container mx-auto px-4 md:px-6 lg:px-8">
			<header className="flex h-20 w-full shrink-0 items-center px-4 md:px-6">
				<div className="ml-auto flex gap-2">
					<SignedIn>
						<UserButton />
					</SignedIn>
				</div>
			</header>
		</div>
	);
}
