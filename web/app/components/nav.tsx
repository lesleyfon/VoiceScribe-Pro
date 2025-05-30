"use client";

import { SignedIn, UserButton, useAuth } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Nav() {
	const { isSignedIn, getToken } = useAuth();
	const router = useRouter();

	useEffect(() => {
		const getData = async () => {
			try {
				const token = await getToken();
				console.log(token);
				const response = await fetch("http://127.0.0.1:8000/user-info", {
					headers: {
						Authorization: `Bearer ${token}`,
						"Content-Type": "application/json",
					},
				});
				if (!response.ok) {
					throw new Error("API BOMBED");
				}
				const data = await response.json();
				console.log({ data });
			} catch (e) {
				throw new Error(e as string);
			}
		};
		getData();
	}, [getToken]);

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
