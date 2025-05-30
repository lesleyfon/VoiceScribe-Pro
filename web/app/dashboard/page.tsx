import Nav from "@/app/components/nav";

export default async function Dashboard() {
	return (
		<>
			<Nav />
			<div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
				<main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
					<header className="flex justify-end items-center p-4 gap-4 h-16">
						Dashboard
					</header>
				</main>
			</div>
		</>
	);
}
