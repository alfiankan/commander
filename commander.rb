class Glow < Formula
  desc "Commander TUI create, run, share and promt snippets and command with Ease"
  homepage "https://alfiankan.github.io/commander-charts"
  version "0.1.0"

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/alfiankan/commander/releases/download/v0.1.0/cmdr-v0.1.0-darwin-amd64.tar.gz"
    sha256 "127a253cbf7b7b8c74018f074f0452aabb0bb29270e63b34782a951fa46f8487"
  end
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/alfiankan/commander/releases/download/v0.1.0/cmdr-v0.1.0-darwin-arm64.tar.gz"
    sha256 "d7f3fc2a1986efa0e9f116d11b1a6e5fae11d2aaa71428593681e63cfbe65113"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/alfiankan/commander/releases/download/v0.1.0/cmdr-v0.1.0-linux-amd64.tar.gz"
    sha256 "97bf347820972cdc247126382be095cf87d99ca4c38391a1c35e897cbcf95297"
  end
  if OS.linux? && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/alfiankan/commander/releases/download/v0.1.0/cmdr-v0.1.0-linux-arm64.tar.gz"
    sha256 "4511e587ce86fced027152391551c38f6feea5a9afb2e3cbe845f465d96b9cb0"
  end

  def install
    bin.install "cmdr"
  end
end
