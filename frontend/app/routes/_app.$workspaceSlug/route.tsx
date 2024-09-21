import { Gauge, FileBox, Settings, LockKeyhole } from "lucide-react";
import { Outlet } from "@remix-run/react";
import React from "react";
import { SiAwsorganizations } from "@icons-pack/react-simple-icons";

export default function Workspace() {
  return (
    <div className="h-full w-full grid grid-rows-1 overflow-hidden">
      <div className="grid grid-cols-[5rem_1fr] overflow-hidden">
        <div className="h-full overflow-hidden hidden md:flex flex-col bg-black border-r border-border">
          <div className="flex flex-col justify-between h-full">
            <div className="flex flex-col px-4 mt-6">
              <a
                href=""
                className="relative p-2 flex items-center md:justify-center border border-border bg-background hover:bg-background hover:border-border hover:border"
              >
                <SiAwsorganizations className="w-6 h-6" />
              </a>
              <nav className="mt-8">
                <ul className="flex flex-col gap-3">
                  <li>
                    <a
                      href=""
                      className="relative p-2 flex items-center md:justify-center border border-border bg-background hover:bg-background hover:border-border hover:border"
                    >
                      <Gauge className="w-6 h-6" />
                    </a>
                  </li>
                  <li>
                    <a
                      href=""
                      className="relative  p-2 flex items-center md:justify-center hover:bg-background hover:border-border hover:border"
                    >
                      <FileBox className="w-6 h-6" />
                    </a>
                  </li>
                  <li>
                    <a
                      href=""
                      className="relative  p-2 flex items-center md:justify-center hover:bg-background hover:border-border hover:border"
                    >
                      <LockKeyhole className="w-6 h-6" />
                    </a>
                  </li>
                  <li>
                    <a
                      href=""
                      className="relative  p-2 flex items-center md:justify-center hover:bg-background hover:border-border hover:border"
                    >
                      <Settings className="w-6 h-6" />
                    </a>
                  </li>
                </ul>
              </nav>
            </div>
            <div className="px-4 mb-4">
              <a
                href=""
                className="relative p-2 flex items-center md:justify-center border border-border bg-background hover:bg-background hover:border-border hover:border"
              >
                <Gauge className="w-6 h-6" />
              </a>
            </div>
          </div>
        </div>
        <div className="grid grid-rows-1 overflow-hidden">
          <div className="grid overflow-hidden grid-rows-[2.5rem_1fr]">
            <div className="w-full overflow-y-auto p-3 scrollbar-thin scrollbar-track-transparent scrollbar-thumb-black/60">
              <Outlet />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
