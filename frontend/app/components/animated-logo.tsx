import { motion } from "framer-motion";
import { cn } from "~/lib/utils";

export default function AnimatedLogo({ className }: { className?: string }) {
  return (
    <motion.div
      className="flex items-center justify-center"
      variants={{
        animate: {
          transition: {
            staggerChildren: 0.3,
          },
        },
      }}
      initial="initial"
      animate="animate"
    >
      <motion.h1
        className={cn("text-md font-bold text-white font-mono", className)}
      >
        {["s", "s", "s"].map((letter, index) => (
          <motion.span
            key={index}
            variants={{
              initial: { y: 0 },
              animate: {
                y: [0, -10, 0],
                transition: {
                  duration: 1.5,
                  ease: "easeInOut",
                  repeatDelay: 0.5,
                  repeat: Infinity,
                  repeatType: "loop",
                },
              },
            }}
            className="inline-block"
            style={{ display: "inline-block" }}
          >
            {letter}
          </motion.span>
        ))}
      </motion.h1>
    </motion.div>
  );
}
