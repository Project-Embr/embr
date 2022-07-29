import React from "react";
import type { NextPage } from "next";
import Head from "next/head";

async function createEmbr() {
  const res = await fetch(`/api/embr`);
}

const Home: NextPage = () => {
  return (
    <div>
      <Head>
        <title>Embr</title>
        <meta
          name="description"
          content="Project Embr is a MicroVPS provider built by UVic VikeLabs students"
        />
        <link rel="icon" href="/favicon.svg" />
      </Head>
      <div className="flex flex-col">
        <h1>Hello World!</h1>
      </div>
      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        onClick={() => createEmbr()}
      >
        Create Embr!
      </button>

      {/* <Main /> */}
    </div>
  );
};

export default Home;
