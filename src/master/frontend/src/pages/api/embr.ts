import type { NextApiRequest, NextApiResponse } from "next";
const { Etcd3 } = require("etcd3");
const client = new Etcd3();

type Data = {
  status: string;
};

/**
 * embr: API handler for the /api/embr route
 * @param {NextApiRequest} req
 * @param {NextApiResponse} res
 */
export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>
) {
  await client.put("/embr/config").value('{"VCpuCount":1,"MemSizeMib":128}');
  res.status(200).json({ status: "success" });
}
